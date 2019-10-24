package nixtui

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rivo/tview"
)

// cmd := exec.Command(
// 	"nix-env"
// 	// -qa --json --show-trace
// )

func searchPackageFilesByName(dirpath string, phrase string) []string {
	contents, err := ioutil.ReadDir(dirpath)
	if err != nil {
		log.Println(err) // FIXME
	}
	if len(contents) == 0 {
		return nil
	}
	result := []string{}
	for _, info := range contents {
		name := info.Name()
		if !info.IsDir() {
			continue
		}
		childPath := filepath.Join(dirpath, name)
		_, defaultErr := os.Stat(filepath.Join(childPath, "default.nix"))
		if defaultErr == nil {
			if strings.Contains(name, phrase) {
				result = append(result, childPath)
			}
		} else {
			result = append(result, searchPackageFilesByName(childPath, phrase)...)
		}
	}
	return result
}

func Main() {
	app := tview.NewApplication()

	result_table := tview.NewTable()
	result_table.SetTitle(" Results ")
	// result_table.SetBorders(true)
	result_table.SetBorder(true)
	result_table.SetSelectable(true, false)
	result_table.SetSelectionChangedFunc(func(row int, col int) {
		_, colOffset := result_table.GetOffset()
		if row > 2 {
			result_table.SetOffset(row-2, colOffset)
		}
	})

	search_phraseInput := tview.NewInputField().SetLabel("Phrase: ")
	// search_phraseInput.SetFieldWidth(20)

	// search_phraseForm := tview.NewForm().AddFormItem(search_phraseInput)

	searchInName := true
	searchInDescription := true

	search_form := tview.NewForm().
		AddFormItem(search_phraseInput).
		AddCheckbox("In name", searchInName, func(checked bool) {
			searchInName = checked
		}).
		AddCheckbox("In description", searchInDescription, func(checked bool) {
			searchInDescription = checked
		}).
		AddButton("Search", func() {
			phrase := search_phraseInput.GetText()
			phrase = strings.TrimSpace(phrase)
			if phrase == "" {
				return
			}
			phrase = strings.ToLower(phrase)
			if searchInName {
				filePathList := searchPackageFilesByName(nixpkgs+"/pkgs", phrase)
				for i, fpath := range filePathList {
					rowI := i + 1
					pathParts := strings.Split(fpath, "/")
					pkgName := pathParts[len(pathParts)-1]
					result_table.SetCell(rowI, 0, tview.NewTableCell(""))
					result_table.SetCell(rowI, 1, tview.NewTableCell(pkgName))
				}
			}
		}).
		AddButton("Quit", func() {
			app.Stop()
		}).
		SetButtonsAlign(tview.AlignRight)
	search_form.SetTitle(" Search ")
	search_form.SetBorder(true)
	search_form.SetTitleAlign(tview.AlignLeft)
	// search_form.SetBorderPadding(0, 0, 0, 0)

	operation_table := tview.NewTable()
	operation_table.SetTitle(" Operations ")
	operation_table.SetBorder(true)

	left_vflex := tview.NewFlex()
	left_vflex.SetDirection(tview.FlexRow)
	left_vflex.AddItem(search_form, 0, 1, true)
	left_vflex.AddItem(operation_table, 0, 1, true)

	main_hflex := tview.NewFlex()
	main_hflex.AddItem(left_vflex, 0, 1, true)
	main_hflex.AddItem(result_table, 0, 3, true)

	result_table.SetCell(0, 0, tview.NewTableCell("").SetExpansion(1))              // status symbol (empty, i, +, a+, d, u, etc)
	result_table.SetCell(0, 1, tview.NewTableCell(" Name").SetExpansion(10))        // package name
	result_table.SetCell(0, 2, tview.NewTableCell(" Installed").SetExpansion(5))    // package installed version
	result_table.SetCell(0, 3, tview.NewTableCell(" Available").SetExpansion(5))    // package available version
	result_table.SetCell(0, 4, tview.NewTableCell(" Description").SetExpansion(30)) // pavkage description

	if err := app.SetRoot(main_hflex, true).Run(); err != nil {
		panic(err)
	}
}
