package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"image/color"
	"strconv"
	"strings"
)

type Client struct {
	Id             uint
	FirstName      string
	Surname        string
	LastName       string
	Remark         string
	Age            int
	Gender         string
	Email          string
	Phone          string
	Address        string
	MembershipType string
}

func main() {
	// Создание окна, настройка размеров, названия и темы.
	a := app.NewWithID("com.example.fitnessapp")
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("Fitness App")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(600, 650))

	ic, _ := fyne.LoadResourceFromPath("icon.png")
	w.SetIcon(ic)

	//СМЕНА ТЕМЫ

	var clients []Client
	var filteredClients []Client
	var createContent *fyne.Container
	var clientsList *widget.List
	var main_bar *fyne.Container
	var settings_bar *fyne.Container

	settings_bar = container.NewVBox(
		container.NewHBox(
			widget.NewButton("Светлая\nтема", func() {
				a.Settings().SetTheme(theme.LightTheme())
			}),
			widget.NewButton("Тёмная\nтема", func() {
				a.Settings().SetTheme(theme.DarkTheme())
			}),
		),
		layout.NewSpacer(),
		widget.NewButtonWithIcon("", theme.MailReplyIcon(), func() {
			w.SetContent(main_bar)
		}),
	)

	settings_btn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		w.SetContent(settings_bar)
	})

	DB, _ := gorm.Open(sqlite.Open("fitness.db"), &gorm.Config{})
	DB.AutoMigrate(&Client{})
	DB.Find(&clients)

	filteredClients = clients

	noClientsLabel := widget.NewLabel("Список клиентов пуст")

	updateClientListLabel(noClientsLabel, len(clients))

	clientsList = widget.NewList(
		func() int {
			return len(filteredClients)
		},
		func() fyne.CanvasObject {
			// Создание контейнера для имени и фамилии
			firstNameLabel := widget.NewLabel("FirstName")
			lastNameLabel := widget.NewLabel("LastName")
			return container.NewHBox(firstNameLabel, lastNameLabel)
		},
		// ОТОБРАЖЕНИЕ ИМЕНИ И ФАМИЛИИ В СПИСКЕ
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			labels := co.(*fyne.Container).Objects
			labels[0].(*widget.Label).SetText(filteredClients[lii].FirstName)
			labels[1].(*widget.Label).SetText(filteredClients[lii].LastName)
		},
	)

	clientsList.OnSelected = func(id widget.ListItemID) {
		delalis_bar := container.NewHBox(
			widget.NewLabel(fmt.Sprintf("Подробнее о клиенте: %s ", filteredClients[id].LastName)),
			layout.NewSpacer(),

			widget.NewButtonWithIcon("Вернуться", theme.ContentUndoIcon(), func() {
				clientsList.UnselectAll()
				w.SetContent(main_bar)
			}),
		)

		// Создаем метки с использованием данных клиента и устанавливаем курсивный стиль текста
		clientFirstName := widget.NewLabel("Имя:  " + filteredClients[id].FirstName)
		clientFirstName.TextStyle = fyne.TextStyle{Italic: true}

		clientSurname := widget.NewLabel("Отчество:  " + filteredClients[id].Surname)
		clientSurname.TextStyle = fyne.TextStyle{Italic: true}

		clientLastName := widget.NewLabel("Фамилия:  " + filteredClients[id].LastName)
		clientLastName.TextStyle = fyne.TextStyle{Italic: true}

		clientRemark := widget.NewLabel("Комментарий:  " + filteredClients[id].Remark)
		clientRemark.TextStyle = fyne.TextStyle{Italic: true}
		clientRemark.Wrapping = fyne.TextWrapBreak

		clientAge := widget.NewLabel(fmt.Sprintf("Возраст: %d ", filteredClients[id].Age))
		clientAge.TextStyle = fyne.TextStyle{Italic: true}

		clientGender := widget.NewLabel("Пол:  " + filteredClients[id].Gender)
		clientGender.TextStyle = fyne.TextStyle{Italic: true}

		clientEmail := widget.NewLabel("Email:  " + filteredClients[id].Email)
		clientEmail.TextStyle = fyne.TextStyle{Italic: true}

		clientPhone := widget.NewLabel("Телефон:  " + filteredClients[id].Phone)
		clientPhone.TextStyle = fyne.TextStyle{Italic: true}

		clientAddress := widget.NewLabel("Адрес:  " + filteredClients[id].Address)
		clientAddress.TextStyle = fyne.TextStyle{Italic: true}

		clientMembershipType := widget.NewLabel("Тип членства:  " + filteredClients[id].MembershipType)
		clientMembershipType.TextStyle = fyne.TextStyle{Italic: true}

		buttonsBox := container.NewHBox(
			// DELETE
			widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				dialog.ShowCustomConfirm(
					"Подтверждение",
					"Да",
					"Нет",
					widget.NewLabel("Вы уверены что хотите удалить запись?"),
					func(b bool) {
						if b {
							DB.Delete(&Client{}, "id = ?", filteredClients[id].Id)
							DB.Find(&clients)

							updateClientListLabel(noClientsLabel, len(clients))
							clientsList.Refresh()

						}
						refreshClientsList(DB, &clients, &filteredClients, clientsList, noClientsLabel)
						w.SetContent(main_bar)
						clientsList.UnselectAll()
					},
					w,
				)
			}),
			// EDIT
			widget.NewButtonWithIcon("", theme.DocumentCreateIcon(),
				func() {
					edit_bar := container.NewHBox(
						widget.NewLabel(fmt.Sprintf("Редактировать данные клиента: %s ", filteredClients[id].LastName)),
						layout.NewSpacer(),

						widget.NewButtonWithIcon("Вернуться", theme.ContentUndoIcon(), func() {
							clientsList.UnselectAll()
							w.SetContent(main_bar)
						}),
					)
					// Создание виджетов для каждого поля структуры Client
					editFirstName := widget.NewEntry()
					editSurname := widget.NewEntry()
					editLastName := widget.NewEntry()
					editRemark := widget.NewMultiLineEntry()
					editAge := widget.NewEntry()
					editGender := widget.NewEntry()
					editEmail := widget.NewEntry()
					editPhone := widget.NewEntry()
					editAddress := widget.NewEntry()
					editMembershipType := widget.NewEntry()

					// Установка значений текстовых полей из экземпляра структуры Client
					editFirstName.SetText(filteredClients[id].FirstName)
					editSurname.SetText(filteredClients[id].Surname)
					editLastName.SetText(filteredClients[id].LastName)
					editRemark.SetText(filteredClients[id].Remark)
					editAge.SetText(strconv.Itoa(filteredClients[id].Age))
					editGender.SetText(filteredClients[id].Gender)
					editEmail.SetText(filteredClients[id].Email)
					editPhone.SetText(filteredClients[id].Phone)
					editAddress.SetText(filteredClients[id].Address)
					editMembershipType.SetText(filteredClients[id].MembershipType)
					// СОХРАНЕНИЕ ОТРЕДАКТИРОВАННОЙ ИНФОРМАЦИИ
					editButton := widget.NewButtonWithIcon(
						"Сохранить изменения",
						theme.ConfirmIcon(),
						func() {
							// Преобразование отредактированного возраста из строки в число
							age, err := strconv.Atoi(editAge.Text)
							if err != nil {
								dialog.ShowError(fmt.Errorf("Возраст должен быть числом"), w)
								return
							}
							DB.Model(&Client{}).Where("id = ?", clients[id].Id).Updates(Client{
								FirstName:      editFirstName.Text,
								Surname:        editSurname.Text,
								LastName:       editLastName.Text,
								Remark:         editRemark.Text,
								Age:            age,
								Gender:         editGender.Text,
								Email:          editEmail.Text,
								Phone:          editPhone.Text,
								Address:        editAddress.Text,
								MembershipType: editMembershipType.Text,
							})

							refreshClientsList(DB, &clients, &filteredClients, clientsList, noClientsLabel)
							w.SetContent(main_bar)
							clientsList.UnselectAll()
						})

					editContent := container.NewVBox(
						edit_bar,
						canvas.NewLine(color.RGBA{0, 200, 122, 255}),
						editFirstName,
						editSurname,
						editLastName,
						editRemark,
						editAge,
						editGender,
						editEmail,
						editPhone,
						editAddress,
						editMembershipType,
						editButton,
					)

					w.SetContent(editContent)
				}),
		)

		// ФИНАЛЬНАЯ ЧАСТЬ

		detalisVbox := container.NewVBox(
			delalis_bar,
			canvas.NewLine(color.RGBA{0, 200, 122, 255}),
			clientFirstName,
			clientSurname,
			clientLastName,
			clientRemark,
			clientAge,
			clientGender,
			clientEmail,
			clientPhone,
			clientAddress,
			clientMembershipType,
			buttonsBox,
		)

		w.SetContent(detalisVbox)
	}

	clientsListScroll := container.NewScroll(clientsList)
	clientsListScroll.SetMinSize(fyne.NewSize(500, 500))

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Поиск клиента...")

	searchEntry.OnChanged = func(query string) {
		filteredClients = filterClients(clients, query)
		updateClientListLabel(noClientsLabel, len(filteredClients))
		clientsList.Refresh()
	}

	// ГЛАВНОЕ МЕНЮ
	main_bar = container.NewVBox(
		widget.NewButtonWithIcon("Добавить клиента", theme.ContentAddIcon(), func() {
			w.SetContent(createContent)
		}),
		canvas.NewLine(color.RGBA{0, 200, 122, 255}),
		searchEntry,
		noClientsLabel,
		clientsListScroll,
		layout.NewSpacer(),
		settings_btn,
	)

	// Поля ввода для нового клиента

	firstNameEntry := widget.NewEntry()
	firstNameEntry.SetPlaceHolder("Имя...")

	surNameEntry := widget.NewEntry()
	surNameEntry.SetPlaceHolder("Отчество...")

	lastNameEntry := widget.NewEntry()
	lastNameEntry.SetPlaceHolder("Фамилия...")

	remarkEntry := widget.NewMultiLineEntry()
	remarkEntry.SetPlaceHolder("Примечания, заболевания и тд...")

	ageEntry := widget.NewEntry()
	ageEntry.SetPlaceHolder("Возраст...")

	genderEntry := widget.NewEntry()
	genderEntry.SetPlaceHolder("Пол...")

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("Электронная почта...")

	phoneEntry := widget.NewEntry()
	phoneEntry.SetPlaceHolder("Телефон...")

	addressEntry := widget.NewEntry()
	addressEntry.SetPlaceHolder("Адрес...")

	membershipEntry := widget.NewEntry()
	membershipEntry.SetPlaceHolder("Дата окончания абонемента...")

	// КНОПКА СОХРАНЕНИЯ ВВЕДЕНЫХ ДАННЫХ
	saveClientBtn := widget.NewButtonWithIcon("Сохранить данные",
		theme.ConfirmIcon(),
		func() {
			age, err := strconv.Atoi(ageEntry.Text)
			if err != nil {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Ошибка",
					Content: "Возраст должен быть числом",
				})
				return
			}
			newClient := Client{
				FirstName:      firstNameEntry.Text,
				Surname:        surNameEntry.Text,
				LastName:       lastNameEntry.Text,
				Remark:         remarkEntry.Text,
				Age:            age,
				Gender:         genderEntry.Text,
				Email:          emailEntry.Text,
				Phone:          phoneEntry.Text,
				Address:        addressEntry.Text,
				MembershipType: membershipEntry.Text,
			}

			DB.Create(&newClient)
			refreshClientsList(DB, &clients, &filteredClients, clientsList, noClientsLabel)
			Refresher(
				firstNameEntry,
				surNameEntry,
				lastNameEntry,
				remarkEntry,
				ageEntry,
				genderEntry,
				emailEntry,
				phoneEntry,
				addressEntry,
				membershipEntry,
			)

			clientsList.UnselectAll()
			w.SetContent(main_bar)
		})

	// Добавление клиента в базу данных
	create_bar := container.NewHBox(
		widget.NewLabel("Добавление клиента"),
		layout.NewSpacer(),
		widget.NewButtonWithIcon("Очистить запись", theme.ContentClearIcon(), func() {
			Refresher(
				firstNameEntry,
				surNameEntry,
				lastNameEntry,
				remarkEntry,
				ageEntry,
				genderEntry,
				emailEntry,
				phoneEntry,
				addressEntry,
				membershipEntry,
			)
		}),
		widget.NewButtonWithIcon("Вернуться", theme.ContentUndoIcon(), func() {
			clientsList.UnselectAll()
			w.SetContent(main_bar)
		}),
	)
	// Добавление клиента в базу данных 2
	createContent = container.NewVBox(
		create_bar,
		canvas.NewLine(color.RGBA{0, 200, 122, 255}),
		container.NewVBox(
			firstNameEntry,
			surNameEntry,
			lastNameEntry,
			remarkEntry,
			ageEntry,
			genderEntry,
			emailEntry,
			phoneEntry,
			addressEntry,
			membershipEntry,
			canvas.NewLine(color.RGBA{0, 200, 122, 255}),
			saveClientBtn,
		),
	)

	w.SetContent(main_bar)

	w.ShowAndRun()
}

func refreshClientsList(DB *gorm.DB, clients *[]Client, filteredClients *[]Client, clientsList *widget.List, noClientsLabel *widget.Label) {
	DB.Find(clients)
	*filteredClients = *clients
	updateClientListLabel(noClientsLabel, len(*filteredClients))
	clientsList.Refresh()
}

func Refresher(entries ...*widget.Entry) {
	for _, entry := range entries {
		entry.SetText("")
		entry.Refresh()
	}
}

func updateClientListLabel(label *widget.Label, count int) {
	if count == 0 {
		label.SetText("Список клиентов пуст")
	} else {
		label.SetText("Всего клиентов: " + strconv.Itoa(count))
	}
	label.Refresh()
}

func filterClients(clients []Client, query string) []Client {
	var filtered []Client
	for _, client := range clients {
		if strings.Contains(strings.ToLower(client.FirstName), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(client.Surname), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(client.LastName), strings.ToLower(query)) {
			filtered = append(filtered, client)
		}
	}
	return filtered
}
