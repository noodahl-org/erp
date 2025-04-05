package mobileapp

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	xw "fyne.io/x/fyne/widget"
	"github.com/noodahl-org/erp/fyne/internal/widgets"
	"github.com/noodahl-org/erp/internal/models"
	"github.com/samber/lo"
)

func (a *MobileApp) ShowEquipmentCatalog() *fyne.Container {
	equip, err := a.erp.FetchEquipment(models.Equipment{})
	if err != nil {
		log.Fatal(err) //todo fix error handling
	}

	sort.Slice(equip, func(i, j int) bool {
		if equip[i].Make != equip[j].Make {
			return equip[i].Make < equip[j].Make
		}
		return equip[i].Model < equip[j].Model
	})

	list := widgets.ToListWidget(equip)

	list.OnSelected = func(id widget.ListItemID) {
		e := equip[id]
		a.ShowEquipmentDetails(&e)
		list.UnselectAll()
	}

	// Add a search field above the list
	options := lo.Map(equip, func(e models.Equipment, _ int) string {
		return fmt.Sprintf("%s %s", e.Model, e.Make)
	})

	search := xw.NewCompletionEntry(options)
	search.SetPlaceHolder("Search equipment...")
	search.OnChanged = func(s string) {
		opts := lo.Filter(options, func(item string, _ int) bool {
			return strings.Contains(strings.ToLower(item), s)
		})

		search.SetOptions(opts)
		search.ShowCompletion()
	}

	// Return the complete container
	return container.NewBorder(
		container.NewVBox(
			container.NewHBox(
				widget.NewLabelWithStyle("Catalog", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				layout.NewSpacer(),
				widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
					a.ShowAddEquipment()
				}),
			),
			search,
		),
		nil, nil, nil,
		container.NewVScroll(list),
	)
}

func (a *MobileApp) ShowEquipmentDetails(equipment *models.Equipment) {
	// Create a detailed view for the equipment
	detailView := container.NewVBox(
		// Header with equipment name
		widget.NewLabelWithStyle("Equipment Details", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		// Grid with equipment details
		container.NewBorder(
			container.NewHBox(
				widget.NewLabel(fmt.Sprintf("%s %s", equipment.Make, equipment.Model)),
				widget.NewLabel(strings.Join(equipment.Tags, ", ")),
				widget.NewLabel(equipment.Category),
			),
			nil, nil, nil,
		),
	)

	// Create a container with the detail view and navigation
	content := container.New(
		layout.NewGridLayoutWithRows(2),
		container.NewVScroll(detailView), // Make details scrollable
		a.Nav,
	)

	// Set the window content to this new view
	a.Main.SetContent(content)
}

func (a *MobileApp) UserEquipmentCatalog() *fyne.Container {
	equip, err := a.erp.FetchUserEquipment(models.UserEquipment{UserID: a.user.ID})
	if err != nil {
		log.Fatal(err) //todo handle error beter
	}
	sort.Slice(equip, func(i, j int) bool {
		if equip[i].Equipment.Make != equip[j].Equipment.Make {
			return equip[i].Equipment.Make < equip[j].Equipment.Make
		}
		return equip[i].Equipment.Model < equip[j].Equipment.Model
	})

	list := widgets.ToListWidget(equip)

	list.OnSelected = func(id widget.ListItemID) {
		e := equip[id]
		a.ShowUserEquipmentDetails(&e)
		list.UnselectAll()
	}

	return container.NewBorder(
		container.NewVBox(
			container.NewHBox(
				widget.NewLabelWithStyle("Equipment Details", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				layout.NewSpacer(),
				widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
					a.ShowUserEquipmentAdd()
				}),
			),
		),
		nil, nil, nil,
		container.NewVScroll(list),
	)
}

func (a *MobileApp) ShowUserEquipmentDetails(equip *models.UserEquipment) {
	sn := ""
	snbinding := binding.BindString(&sn)
	snbinding.Set(equip.SerialNumber)

	yr := ""
	yrbinding := binding.BindString(&yr)
	yrbinding.Set(fmt.Sprintf("%v", equip.Year))

	content := container.NewBorder(
		container.NewVBox(
			container.NewHBox(
				widget.NewLabelWithStyle("Equipment Details", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				layout.NewSpacer(),
				widget.NewButtonWithIcon("Maint Log", theme.ContentAddIcon(), func() {}),
				widget.NewButtonWithIcon("Delete", theme.ContentAddIcon(), func() {
					if err := a.erp.DeleteUserEquipment(equip.ID); err != nil {
						log.Fatal(err) //todo handle errs

					}
					//todo fix this refactor user equipment catalog
					//into its own show
					a.Main.SetContent(
						container.NewGridWithRows(2,
							a.UserEquipmentCatalog(),
							a.Nav,
						),
					)
				}),
			),

			container.NewBorder(nil, nil,
				container.NewGridWithColumns(2,
					widget.NewLabelWithStyle("Make", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
					widget.NewLabel(equip.Equipment.Make),
					widget.NewLabelWithStyle("Model", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
					widget.NewLabel(equip.Equipment.Model),
					widget.NewLabelWithStyle("Added", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
					widget.NewLabel(equip.CreatedAt.Format("2006-02-01")),
				),
				container.NewGridWithColumns(2,
					widget.NewLabelWithStyle("Category", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
					widget.NewLabel(equip.Equipment.Category),
					widget.NewLabelWithStyle("Tags", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
					widget.NewLabel(strings.Join(equip.Equipment.Tags, ", ")),
				),
			),

			widget.NewForm(
				[]*widget.FormItem{
					{Text: "serial num", Widget: widget.NewEntryWithData(snbinding)},
					{Text: "year", Widget: widget.NewEntryWithData(yrbinding)},
					{Text: "", Widget: widget.NewButton("Save", func() {
						yearstr, _ := yrbinding.Get()
						year, err := strconv.Atoi(yearstr)
						if err != nil {
							log.Fatal(err)
						}
						equip.SerialNumber, _ = snbinding.Get()
						equip.Year = year

						equip, err = a.erp.UpserUserEquipment(*equip)
						if err != nil {
							log.Fatal(err)
						}
					})},
				}...,
			),
		),
		nil, nil, nil,
	)

	a.Main.SetContent(container.New(
		layout.NewGridLayoutWithRows(2),
		content,
		a.Nav,
	))
}

func (a *MobileApp) ShowAddEquipment() {
	equip := models.Equipment{}
	makebinding := binding.BindString(&equip.Make)
	modelbinding := binding.BindString(&equip.Model)
	categorybinding := binding.BindString(&equip.Category)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Equipment Details", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewForm(
				[]*widget.FormItem{
					{Text: "Make", Widget: widget.NewEntryWithData(makebinding)},
					{Text: "Model", Widget: widget.NewEntryWithData(modelbinding)},
					{Text: "Category", Widget: widget.NewEntryWithData(categorybinding)},
					{Text: "", Widget: widget.NewButtonWithIcon("Save", theme.ConfirmIcon(), func() {
						result, err := a.erp.UpsertEquipment(equip)
						if err != nil {
							log.Fatal(err) //todo handle errs
						}
						a.ShowEquipmentDetails(result)
					})},
					{Text: "", Widget: widget.NewButtonWithIcon("Cancel", theme.ContentUndoIcon(), func() {
						//todo refactor
						a.Main.SetContent(
							container.NewGridWithRows(2,
								a.ShowEquipmentCatalog(),
								a.Nav,
							),
						)
					})},
				}...,
			),
		), nil, nil, nil,
	)

	a.Main.SetContent(container.New(
		layout.NewGridLayoutWithRows(2),
		content,
		a.Nav,
	))
}

func (a *MobileApp) ShowUserEquipmentAdd() {
	userEquip := models.UserEquipment{
		UserID: a.user.ID,
	}
	equiplabel := ""
	equipbinding := binding.BindString(&equiplabel)

	equipment, err := a.erp.FetchEquipment(models.Equipment{})
	if err != nil {
		log.Fatal(err) //todo fix error handling
	}

	sort.Slice(equipment, func(i, j int) bool {
		if equipment[i].Make != equipment[j].Make {
			return equipment[i].Make < equipment[j].Make
		}
		return equipment[i].Model < equipment[j].Model
	})

	// Add a search field above the list
	options := lo.Map(equipment, func(e models.Equipment, _ int) string {
		return fmt.Sprintf("%s %s", e.Model, e.Make)
	})

	search := xw.NewCompletionEntry(options)
	search.SetPlaceHolder("Search equipment...")
	search.OnChanged = func(s string) {
		opts := lo.Filter(options, func(item string, _ int) bool {
			return strings.Contains(strings.ToLower(item), s)
		})
		if len(opts) == 1 {
			if equip, _, ok := lo.FindIndexOf(equipment, func(e models.Equipment) bool {
				return fmt.Sprintf("%s %s", e.Model, e.Make) == opts[0]
			}); ok {
				userEquip.EquipmentID = equip.ID
				userEquip.Equipment = equip
				equipbinding.Set(fmt.Sprintf("%s %s", equip.Model, equip.Make))
			}
		}

		search.SetOptions(opts)
		search.ShowCompletion()
	}

	content := container.NewBorder(
		container.NewVBox(
			container.NewBorder(
				container.NewGridWithRows(2,
					widget.NewLabelWithStyle("Add to Inventory", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
					search,
				), nil, nil, nil,
			),

			widget.NewLabelWithStyle("Inventory Details", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewForm(
				[]*widget.FormItem{
					{Text: "Equipment", Widget: widget.NewLabelWithData(equipbinding)},
					// {Text: "Model", Widget: widget.NewEntryWithData(modelbinding)},
					// {Text: "Category", Widget: widget.NewEntryWithData(categorybinding)},
					{Text: "", Widget: widget.NewButtonWithIcon("Save", theme.ConfirmIcon(), func() {
						result, err := a.erp.UpserUserEquipment(userEquip)
						if err != nil {
							log.Fatal(err) //todo handle errs
						}
						a.ShowUserEquipmentDetails(result)
					})},
					{Text: "", Widget: widget.NewButtonWithIcon("Cancel", theme.ContentUndoIcon(), func() {
						//todo refactor
						a.Main.SetContent(
							container.NewGridWithRows(2,
								a.ShowEquipmentCatalog(),
								a.Nav,
							),
						)
					})},
				}...,
			),
		), nil, nil, nil,
	)

	a.Main.SetContent(container.New(
		layout.NewGridLayoutWithRows(2),
		content,
		a.Nav,
	))
}
