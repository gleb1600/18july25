# ZipArchive: Сервис для упаковки файлов из открытых источников в zip-архив.

## Особенности сервиса:
- Возможность создавать zip-архивы и добавлять в них задачи (добавлять файлы из открытого доступа в zip-архив)
- В работе может находиться не боллее 3 zip-архивов одновременно
- В каждом zip-архиве могут находить не более 3 задач
- Возможность получения статуса каждого zip-архива и каждой задачи
- Упраление производиться обращением к localhost
- После архивации: получение ссылки на zip-архив

## Работа

**1.** **Запуск**:
```bash
go run main.go
```
**2.** **Создание нового zip-архива**:
```bash
http://localhost:8080/createziparchive
```
*Ответ: Возвращает ID и статус zip-архива:*
```bash
ZipArchiveID: 1755555550000000000, Status: ZA Created
```
**3.** **Просмотр всех имеющихся zip-архивов**:
```bash
http://localhost:8080/ziparchives
```
*Ответ: Возвращает ID, статус и количество задач для каждого zip-архива:*
```bash
ZipArchiveID: 1755555550000000000, Status: ZA Created, TasksNumber: 0
ZipArchiveID: 1755555551000000000, Status: ZA Created, TasksNumber: 0
ZipArchiveID: 1755555552000000000, Status: ZA Created, TasksNumber: 0
```
**4.** **Создание нового задания для zip-архива**:

*Выполняется POST-запросом, в теле передается ссылка на интересующий файл*
```bash
curl -X POST -d "url=https://pictures/tiger.jpeg" http://localhost:8080/createtask/<ZipArchiveID>
```
*Например для zip-архива c ID:1755555551000000000*
```bash
curl -X POST -d "url=https://pictures/tiger.jpeg" http://localhost:8080/createtask/1755555551000000000
```
*Ответ: Возвращает ID родительского zip-архива, ID самого задания и статус:*
```bash
ZipArchiveID: 1755555551000000000, TaskID: 1755555553000000000, Status: Task Created
```
**5.** **Просмотр статуса задания**:
```bash
http://localhost:8080/task/<TaskID>
```
*Например для задания c ID:1755555553000000000*
```bash
http://localhost:8080/task/1755555553000000000
```
*Ответ: Возвращает ID задания и статус:*
```bash
TaskID: 1755555553000000000, Status: Task Completed Successfully
```
**6.** **Просмотр статуса zip-архива**:
```bash
http://localhost:8080/ziparchive/<ZipArchiveID>
```
*Например для zip-архива c ID:1755555551000000000*
```bash
http://localhost:8080/ziparchive/1755555551000000000
```
*Ответ: Возвращает ID, статус и количество задач для zip-архива, а также ID и статус заданий:*
```bash
ZipArchiveID: 1755555551000000000, Status: ZA Created, TasksNumber: 2
TaskID: 1755555553000000000, Status: Task Completed Successfully
TaskID: 1755555554000000000, Status: Task Completed Unsuccessfully
```
*При наличии в zip-архиве 3 заданий также возвращается ссылка на готовый zip-архив*
```bash
ZipArchiveID: 1755555551000000000, Status: ZA Created, TasksNumber: 2
TaskID: 1755555553000000000, Status: Task Completed Successfully
TaskID: 1755555554000000000, Status: Task Completed Unsuccessfully
TaskID: 1755555555000000000, Status: Task Completed Successfully
To DOWNLOAD ZipArchive: http://localhost:8080/download/1755555551000000000
```
*Перейдя по download-ссылке скачается готовый zip-архив*

**7.** **Контроль всех имеющихся zip-архивов**:
```bash
http://localhost:8080/ziparchives
```
*Ответ: Возвращает ID, статус и количество задач для каждого zip-архива:*
```bash
ZipArchiveID: 1755555550000000000, Status: ZA Created, TasksNumber: 0
ZipArchiveID: 1755555551000000000, Status: ZA Completed Successfully, TasksNumber: 3
ZipArchiveID: 1755555552000000000, Status: ZA Created, TasksNumber: 0
ZipArchiveID: 1755555556000000000, Status: ZA Created, TasksNumber: 2
ZipArchiveID: 1755555557000000000, Status: ZA Completed Successfully, TasksNumber: 3
```
*Одновременно в работе (не завершены) могу быть не более 3*
## Структура проекта
```bash
ziparchive
        ├── handlers
        │   └── handlers.go               # HTTP-обработчики
        ├── storage
        │   ├── task.go                   # Реализация структуры Задания
        │   ├── zipArchiveManager.go      # Реализация структуры Менеджера Zip-архивов
        │   └── zipArchive.go             # Реализация структуры Zip-архива
        │
        ├── main.go                       # Точка входа
        └── go.mod
```
