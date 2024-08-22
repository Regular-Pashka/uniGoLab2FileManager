package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func createFile(path string) {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println("1. Создать каталог")
    fmt.Println("2. Создать файл")
    scanner.Scan()
    ans := scanner.Text()
    switch ans {
    case "1":
        err := os.MkdirAll(path, os.ModePerm)
        if err != nil {
            fmt.Println("Каталог не был создан:", err)
        } else {
            fmt.Println("Каталог успешно создан")
        }
    case "2":
        parentDir := filepath.Dir(path)
        if parentDir != "" {
            err := os.MkdirAll(parentDir, os.ModePerm)
            if err != nil {
                fmt.Println("Каталог не был создан:", err)
                return
            }
        }
        _, err := os.Create(path)
        if err != nil {
            fmt.Println("Файл не был создан:", err)
        } else {
            fmt.Println("Файл успешно создан")
        }
    default:
        fmt.Println("Неверная команда!\nФайл/каталог не будет создан")
    }
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("не удалось открыть исходный файл: %w", err)
	}
	defer sourceFile.Close()

	fileInfo, err := os.Stat(src)
	if err != nil {
		fmt.Println("Ошибка при получении информации о файле:", err)
		return nil
	}

	destinationFile, err := os.Create(dst + fileInfo.Name())
	if err != nil {
		return fmt.Errorf("не удалось создать файл назначения: %w", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("не удалось скопировать файл: %w", err)
	}

	return nil
}

func listFilesWithExt(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var listPaths []string

	fmt.Print("Введите расширение: ")
	var ext string
	_, err = fmt.Scanln(&ext)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), strings.ToLower(ext)) {
			listPaths = append(listPaths, filepath.Join(path, file.Name()))
		}
	}

	for _, p := range listPaths {
		fmt.Println(filepath.Base(p))
	}

	return nil
}

func deleteFile(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("File doesn't exists: %v", err)
	}

	if !fileInfo.IsDir() {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	} else {
		files, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			err = os.Remove(path)
			if err != nil {
				return err
			}
			return nil
		}

		for i := len(files) - 1; i >= 0; i-- {
			err = deleteFile(filepath.Join(path, files[i].Name()))
			if err != nil {
				return err
			}
		}

		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func findFileInDir(file, dir string) error {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("ошибка доступа к каталогу: %w", err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("%s не является каталогом", dir)
	}

	fileName := filepath.Base(file)
	var foundPaths []string

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if path != dir {
				return filepath.SkipDir
			}
			return nil
		}

		if info.Name() == fileName {
			foundPaths = append(foundPaths, path)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("ошибка при поиске файла: %w", err)
	}

	if len(foundPaths) == 0 {
		fmt.Println("Такого файла нет в каталоге")
	} else {
		for _, path := range foundPaths {
			fmt.Println("Путь к найденному файлу:", path)
		}
	}

	return nil
}

func getPath() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите путь к файлу/каталогу: ")
	path, _ := reader.ReadString('\n')
	return strings.TrimSpace(path)
}

func getFileName() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите имя искомого файла: ")
	name, _ := reader.ReadString('\n')
	return strings.TrimSpace(name)
}

func printMenu() {
	fmt.Println("1. Вывести информацию о файле/каталоге")
	fmt.Println("2. Изменение имени файла/каталога")
	fmt.Println("3. Создание нового файла или каталога по заданному пути")
	fmt.Println("4. Создание копии файла по заданному пути")
	fmt.Println("5. Вывод списка файлов каталога")
	fmt.Println("6. Вывод списка файлов каталога, имеющих определенное расширение")
	fmt.Println("7. Удаление файла или каталога")
	fmt.Println("8. Поиск файла или каталога в выбранном каталоге")
	fmt.Println("9. Выход из программы")
	fmt.Print("> ")
}

func printFileInfo(path string) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("Ошибка при получении информации о файле:", err)
		return
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Ошибка при получении абсолютного пути:", err)
		return
	}
	fmt.Println("Имя:", fileInfo.Name())
	fmt.Println("Путь:", path)
	fmt.Println("Абсолютный путь:", absPath)
	fmt.Println("Время последней модификации", fileInfo.ModTime())
	fmt.Println("Размер файла в байтах:", fileInfo.Size())
	fmt.Println("Режим доступа:", fileInfo.Mode())
	fmt.Println("Время создания:", fileInfo.ModTime())
	fmt.Println("Это каталог?", fileInfo.IsDir())

}

func listFiles(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("ошибка доступа к каталогу: %w", err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("%s не является каталогом", path)
	}

	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("ошибка чтения содержимого каталога: %w", err)
	}

	fmt.Println("Файлы в", path+":")
	for _, fileInfo := range fileInfos {
		fmt.Println(fileInfo.Name())
	}

	return nil
}

func changeName(path string) {
	fmt.Print("Введите новое имя файла: ")
	reader := bufio.NewReader(os.Stdin)
	newName, _ := reader.ReadString('\n')
	newName = strings.TrimSpace(newName)

	err := os.Rename(path, filepath.Dir(path)+"/"+newName)
	if err != nil {
		fmt.Println("Файл не был изменен:", err)
	} else {
		fmt.Println("Файл успешно изменен")
	}
}


func main() {
	for {
		printMenu()

		var ans int
		fmt.Scanln(&ans)

		switch ans {
		case 1:
			path := getPath()
			printFileInfo(path)
		case 2:
			path := getPath()
			changeName(path)
		case 3:
			path := getPath()
			createFile(path)
		case 4:
			path := getPath()
			// copyFile(path, getPath())
			err := copyFile(path, getPath())
			if err != nil {
				fmt.Println("Произошла ошибка при копировании файла:", err)
			}
		case 5:
			path := getPath()
			listFiles(path)
		case 6:
			path := getPath()
			listFilesWithExt(path)
		case 7:
			path := getPath()
			errDel := deleteFile(path)
			if errDel != nil {
				fmt.Println("Ошибка:", errDel)
			}
		case 8:
			path := getFileName()
			fmt.Println("Путь откуда:")
			dir := getPath()
			findFileInDir(path, dir)	
		case 9:
			fmt.Println("Выход из программы")
			return
		default:
			fmt.Println("Неверная команда!")
		}
	}
}


/* 

*/

/* 
TODO:

Практическая работа №2. Работа с файлами.
Задание.
Реализовать простейший файловый менеджер. (Для Java использовать пакеты
java.io|java.nio, класс File|Files). Программа должна обладать следующим
функционалом:
1 Выбор файла или каталога для работы;
2 Вывод абсолютного пути для текущего файла или каталога;
3 Вывод содержимого каталога;
4 Вывод всей возможной информации для заданного файла;
5 Изменение имени файла или каталога;
6 Создание нового файла или каталога по заданному пути;
7 Создание копии файла по заданному пути;
8 Вывод списка файлов текущего каталога, имеющих расширение,
задаваемое пользователем;
9 Удаление файла или каталога;
10.Поиск файла или каталога в выбранном каталоге;
Результаты лабораторной работы оформить в виде отчета с результатами
работы программы.

*/