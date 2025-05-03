package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var Current_User = User{}
var Current_Chat_String string

func main() {
	myApp := app.New()
	Screen := myApp.NewWindow("Hello")

	name_Text := widget.NewLabel("")

	mainUiInnerWindows := container.NewHBox()

	newRoomNameEntry := widget.NewEntry()

	chats := container.NewVBox()
	chat_loaded := false

	rightSideContainer := container.NewVBox()

	messeges := container.NewVBox()
	sendMessegeEntry := widget.NewEntry()
	addUserEntry := widget.NewEntry()

	go func() {
		for true {
			time.Sleep(time.Second)
			if chat_loaded {
				getRoomData, err := json.Marshal(NetworkChat{Current_User.Username, []string{}})

				resp, err := http.Post("http://localhost:5151/GetChatsForUser", "json", bytes.NewBuffer(getRoomData))
				if err != nil {
					panic(err)
				}

				temp_string_chats, err := io.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}

				chats_data := []Chat{}
				if err := json.Unmarshal(temp_string_chats, &chats_data); err != nil {
					panic(err)
				}

				// current_chat := &Chat{}
				current_id := 0

				for ci := range chats_data {
					if chats_data[ci].Name == Current_Chat_String {
						// current_chat = &chats_data[ci]
						current_id = ci
					}
				}

				messeges.RemoveAll()
				for _, m := range chats_data[current_id].Messeges {
					messeges.Add(widget.NewLabel(m.Username + ": " + m.Messege))
				}

				fmt.Println(chats_data)
			}
		}
	}()

	sideBar := container.NewHSplit(
		container.NewVBox(
			name_Text,

			chats,

			widget.NewButton("New Chat", func() {
				mainUiInnerWindows.Add(container.NewInnerWindow("pls",
					container.NewVBox(
						container.NewGridWithColumns(
							2,
							widget.NewLabel("Chat Name"),
							newRoomNameEntry,
						),
						widget.NewButton("Create Chat", func() {
							data := NetworkChat{
								newRoomNameEntry.Text,
								[]string{Current_User.Username},
							}
							temp_data, err := json.Marshal(data)
							if err != nil {
								panic(err)
							}
							_, err = http.Post("http://localhost:5151/MakeChat", "json", bytes.NewBuffer(temp_data))
							if err != nil {
								panic(err)
							}

							getRoomData, err := json.Marshal(NetworkChat{Current_User.Username, []string{}})

							resp, err := http.Post("http://localhost:5151/GetChatsForUser", "json", bytes.NewBuffer(getRoomData))
							if err != nil {
								panic(err)
							}

							temp_string_chats, err := io.ReadAll(resp.Body)
							if err != nil {
								panic(err)
							}

							chats_data := []Chat{}
							if err := json.Unmarshal(temp_string_chats, &chats_data); err != nil {
								panic(err)
							}

							current_chat := &Chat{}

							for ci := range chats_data {
								if chats_data[ci].Name == data.Name {
									current_chat = &chats_data[ci]
								}
							}

							chats.Add(widget.NewButton(data.Name, func() {
								chat_loaded = true
								Current_Chat_String = data.Name
								rightSideContainer.RemoveAll()
								rightSideContainer.Add(
									container.NewVBox(
										container.NewGridWithColumns(3,
											widget.NewLabel(""),
											widget.NewLabel(data.Name),
											widget.NewLabel(""),
										),
										messeges,
										container.NewGridWithColumns(
											2,
											sendMessegeEntry,
											widget.NewButton("send", func() {
												current_chat.Messeges = append(current_chat.Messeges, NewMessege(Current_User.Username, sendMessegeEntry.Text))

												messeges.RemoveAll()
												for _, m := range current_chat.Messeges {
													messeges.Add(widget.NewLabel(m.Username + ": " + m.Messege))
												}

												messege_network_data := NetworkChat{current_chat.Name, []string{Current_User.Username, sendMessegeEntry.Text}}

												newMessegeData, err := json.Marshal(messege_network_data)
												if err != nil {
													panic(err)
												}

												http.Post("http://localhost:5151/SendMessege", "json", bytes.NewBuffer(newMessegeData))
											}),
										),
									),
								)
							}))
						}),
					)))
			}),
		),
		rightSideContainer,
	)
	sideBar.Offset = -0.1

	mainUi := container.NewVBox(
		sideBar,
		mainUiInnerWindows,
	)

	loginUsernameInput := widget.NewEntry()
	loginUsernameInput.SetPlaceHolder("")

	loginUsernameContainer := container.NewGridWithColumns(
		2,
		widget.NewLabel("Username: "),
		loginUsernameInput,
	)

	loginPasswordInput := widget.NewEntry()
	loginPasswordInput.SetPlaceHolder("")

	loginPasswordContainer := container.NewGridWithColumns(2,
		widget.NewLabel("Password: "),
		loginPasswordInput,
	)

	loginPasswordContainer.Resize(fyne.NewSize(100, 100))
	loginUsernameContainer.Resize(fyne.NewSize(100, 100))

	loginUi := container.NewVBox(
		widget.NewLabel("Login"),
		loginUsernameContainer,
		loginPasswordContainer,
		widget.NewButton("Create User", func() {
			data, err := json.Marshal(NewUser(loginUsernameInput.Text, loginPasswordInput.Text))
			if err != nil {
				panic(err)
			}
			resp, err := http.Post("http://localhost:5151/AddUser", "json", bytes.NewBuffer(data))
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			Screen.SetContent(mainUi)

			Current_User = NewUser(loginUsernameInput.Text, loginPasswordInput.Text)
		}),

		widget.NewButton("Login", func() {
			data, err := json.Marshal(NewUser(loginUsernameInput.Text, loginPasswordInput.Text))
			if err != nil {
				panic(err)
			}
			resp, err := http.Post("http://localhost:5151/GetUser", "json", bytes.NewBuffer(data))
			if err != nil {
				panic(err)
			}

			if resp.StatusCode == 202 {
				temp_data, err := io.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}

				temp_user := User{}

				json.Unmarshal(temp_data, &temp_user)

				Current_User = temp_user
				name_Text.Text += Current_User.Username

				Screen.SetContent(mainUi)

			}

			getRoomData, err := json.Marshal(NetworkChat{Current_User.Username, []string{}})

			resp, err = http.Post("http://localhost:5151/GetChatsForUser", "json", bytes.NewBuffer(getRoomData))
			if err != nil {
				panic(err)
			}

			temp_string_chats, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			chats_data := []Chat{}
			if err := json.Unmarshal(temp_string_chats, &chats_data); err != nil {
				panic(err)
			}

			for _, chat := range chats_data {
				chats.Add(widget.NewButton(chat.Name, func() {
					chat_loaded = true
					messeges.RemoveAll()
					for _, m := range chat.Messeges {
						messeges.Add(widget.NewLabel(m.Username + ": " + m.Messege))
					}

					rightSideContainer.RemoveAll()
					rightSideContainer.Add(
						container.NewVBox(
							container.NewGridWithColumns(3,
								widget.NewLabel(""),
								widget.NewLabel(chat.Name),
								widget.NewLabel(""),
							),
							messeges,
							container.NewGridWithColumns(
								3,
								sendMessegeEntry,
								widget.NewButton("send", func() {
									chat.Messeges = append(chat.Messeges, NewMessege(Current_User.Username, sendMessegeEntry.Text))

									messeges.RemoveAll()
									for _, m := range chat.Messeges {
										messeges.Add(widget.NewLabel(m.Username + ": " + m.Messege))
									}

									messege_network_data := NetworkChat{chat.Name, []string{Current_User.Username, sendMessegeEntry.Text}}

									newMessegeData, err := json.Marshal(messege_network_data)
									if err != nil {
										panic(err)
									}

									http.Post("http://localhost:5151/SendMessege", "json", bytes.NewBuffer(newMessegeData))
								}),
								widget.NewButton("Add User", func() {
									mainUiInnerWindows.Add(container.NewInnerWindow("Add User", container.NewHBox(
										addUserEntry,
										widget.NewButton("Add", func() {
											network_data := NetworkChat{chat.Name, []string{addUserEntry.Text}}

											newMessegeData, err := json.Marshal(network_data)
											if err != nil {
												panic(err)
											}

											http.Post("http://localhost:5151/AddUserToChat", "json", bytes.NewBuffer(newMessegeData))
										}),
									)))
								}),
							),
						),
					)
				}))
			}
		}),
	)

	Screen.SetContent(
		loginUi,
	)

	Screen.Show()
	myApp.Run()
	tidyUp()
}

func tidyUp() {
	fmt.Println("Exited")
}
