# test_task_02_02

Тестовое задание в YADRO, полученное 26.01.2023 с дедлайном 02.02.2023.

Консольная программа, которая читает CSV-файл (comma-separated values) с заголовком, в котором перечислены названия столбцов. Строки нумеруются целыми положительными числами, необязательно в порядке возрастания. В ячейках CSV-файла могут хранится или целые числа, или выражения вида [=ARG1 OP ARG2], где ARG1 и ARG2 – целые числа или адреса ячеек в формате Имя_колонки Номер_строки, а OP – арифметическая операция из списка: +, -, *, /. Программа вычисляет значения ячеек в формате выражений и выводит получившуюся табличку в виде CSV-представления в консоль.

### Запуск
Для запуска программы должен быть установлен компилятор Go. 
Если он установлен, скачайте и распакуйте репозиторий в произвольную директорию, а затем запустите в этой директории терминал. 

Скачать репозиторий (пропишите команду в Git):
```
git clone https://github.com/rabsiyanin/test_task_02_02.git
```

Для компиляции и сборки пропишите команду:
```
go build csvreader.go 
``` 
Для запуска на Windows пропишите команду:
```
csvreader.exe file.csv
```
Для запуска на Linux пропишите команду:
```
./csvreader file.csv
```

### Важные комментарии

* В качестве вводного .csv файла используется file.csv, который можно видоизменять для тестирования приложения. В репозитории приложены и другие .csv файлы с разными сценариями (как положительными, так и отрицательными) представления .csv-таблиц для полноценного тестирования их обработки. Для запуска программы со своим .csv файлом просто вставьте его в директорию и измените аргумент названия файла в вышеупомянутых командах на название интересующего файла.

* В этой программе для обработки ошибок panic() используется гораздо чаще, чем возврат error. Это связано с решением никогда не позволять пользователю продолжать выполнение, когда ввод не является корректным. Тестирование приложения благодаря этому подходу было проще, однако такой подход вреден, когда речь идет о крупномасштабных приложениях, где следует избегать фатального завершения программы. В связи с этим некоторые из случаев обработки ошибок были оставлены с возвращениями error.
