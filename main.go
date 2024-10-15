package main

import (
	storage "DBTask/Storage"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	//Введите свои данные
	user := ""
	pwd := ""
	host := ""
	port := ""
	db := ""
	st, err := storage.New("postgres://" + user + ":" + pwd + "@" + host + ":" + port + "/" + db + "")
	if err != nil {
		fmt.Println("Ошибка подключения")
		return
	}

	var task storage.Task
	fmt.Println("СОЗДАНИЕ ЗАДАЧИ")
	fmt.Println("Введите id автора")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	task.AuthorID, err = strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Ошибка парсинга" + err.Error())
		return
	}
	fmt.Println("Введите id ответственного")
	scanner.Scan()
	task.AssignedID, err = strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Ошибка парсинга" + err.Error())
		return
	}
	fmt.Println("Введите заголовок")
	scanner.Scan()
	task.Title = scanner.Text()
	fmt.Println("Введите описание")
	scanner.Scan()
	task.Content = scanner.Text()
	err = st.NewTask(task)
	if err != nil {
		fmt.Println("Ошибка при создании задачи")
		return
	}

	fmt.Println("ОБНОВЛЕНИЕ ЗАДАЧИ")
	task = storage.Task{}
	end := false
	fmt.Println("Введите id задачи")
	scanner.Scan()
	task.ID, err = strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Ошибка парсинга" + err.Error())
		return
	}
	fmt.Println("Завершить задачу?")
	scanner.Scan()
	end, err = strconv.ParseBool(scanner.Text())
	if err != nil {
		fmt.Println("Ошибка парсинга" + err.Error())
		return
	}
	if end {
		task.Closed = int64(time.Now().Unix())
	} else {
		task.Closed = 0
	}
	err = st.UpdateTask(task.ID, task)
	if err != nil {
		fmt.Println("Ошибка при обновлении задачи")
		return
	}

	fmt.Println("УДАЛЕНИЕ ЗАДАЧИ")
	task = storage.Task{}
	fmt.Println("Введите id задачи")
	scanner.Scan()
	task.ID, err = strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Ошибка парсинга" + err.Error())
		return
	}
	err = st.DeleteTask(task.ID)
	if err != nil {
		fmt.Println("Ошибка при удалении задачи")
		return
	}

	fmt.Println("ВСЕ ЗАДАЧИ")
	allTasks, err := st.AllTasks()
	if err != nil {
		fmt.Println("Ошибка при выборке всех задач")
		return
	}
	fmt.Println(allTasks)

	fmt.Println("ПОИСК ПО АВТОРУ")
	fmt.Println("Введите id автора")
	scanner.Scan()
	task.AuthorID, err = strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Ошибка парсинга" + err.Error())
		return
	}
	allTasks, err = st.TasksOfAuthor(task.AuthorID)
	if err != nil {
		fmt.Println("Ошибка при выборке задач по автору")
		return
	}
	fmt.Println(allTasks)

	fmt.Println("ПОИСК ПО МЕТКЕ")
	fmt.Println("Введите метку")
	scanner.Scan()
	label := scanner.Text()
	allTasks, err = st.TasksOfLabel(label)
	if err != nil {
		fmt.Println("Ошибка при выборке задач по метке")
		return
	}
	fmt.Println(allTasks)

	return
}
