package upgrade

import (
	"github.com/nbkit/mdf/db"
	"github.com/nbkit/mdf/framework/files"
	"github.com/nbkit/mdf/framework/glog"
	"github.com/nbkit/mdf/framework/md"
	"github.com/nbkit/mdf/utils"
	"io/ioutil"
	"os"
	"strings"
)

type ScriptOption struct {
	Path       string
	Extensions []string
}

type FileScript interface {
	Exec()
}

func Script(option ScriptOption) FileScript {
	return &scriptImpl{option: option}
}

type scriptImpl struct {
	option ScriptOption
}

func newExcelSeed() *scriptImpl {
	ser := &scriptImpl{}
	return ser
}
func (s *scriptImpl) Exec() {
	fileList := getAllFiles(utils.JoinCurrentPath(s.option.Path), s.option.Extensions)
	if len(fileList) == 0 {
		return
	}
	for _, f := range fileList {
		if !utils.PathExists(f) {
			continue
		}
		if s.isScript(f) {
			s.runScript(f)
			continue
		}
		if s.isExcel(f) {
			s.runExcel(f)
			continue
		}
	}
}
func (s *scriptImpl) isScript(file string) bool {
	parts := strings.Split(file, ".")
	if strings.ToLower(parts[len(parts)-1]) == "sql" {
		return true
	}
	return false
}

func (s *scriptImpl) isExcel(file string) bool {
	parts := strings.Split(file, ".")
	if strings.ToLower(parts[len(parts)-1]) == "xlsx" {
		return true
	}
	return false
}
func (s *scriptImpl) runScript(file string) error {
	if bt, err := ioutil.ReadFile(file); err != nil {
		return glog.Error(err)
	} else {
		if err := db.Default().Exec(string(bt)).Error; err != nil {
			return glog.Error(err)
		}
	}
	return nil
}
func (s *scriptImpl) runExcel(file string) error {
	if data, err := files.NewExcelSv().GetExcelDatas(file); err != nil {
		return err
	} else {
		return s.handExcelData(data)
	}
}
func (s *scriptImpl) handExcelData(data []files.ImportData) error {
	groups := make(map[string][]files.ImportData)
	for i := range data {
		d := data[i]
		parts := strings.Split(d.SheetName, ".")
		groups[parts[0]] = append(groups[parts[0]], d)
	}
	for k, d := range groups {
		c := utils.NewFlowContext()
		c.Request.Data = d
		c.Request.Action = "import"
		c.Request.Widget = k
		md.ActionSv().DoAction(c)
	}
	return nil
}

func getAllFiles(dirPth string, extensions []string) (files []string) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if !utils.StringIsAlphanumeric(fi.Name()[:1]) {
			continue
		}
		if fi.IsDir() {
			// 读取子目录下文件
			if tempFields := getAllFiles(dirPth+PthSep+fi.Name(), extensions); tempFields != nil && len(tempFields) > 0 {
				files = append(files, tempFields...)
			}
		} else {
			// 过滤指定格式
			parts := strings.Split(fi.Name(), ".")
			ext := strings.ToLower(parts[len(parts)-1])

			if utils.StringsContains([]string{"bak", "back", "log"}, ext) >= 0 {
				continue
			}
			if len(extensions) == 0 || utils.StringsContains(extensions, ext) >= 0 {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}
	return files
}
