# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m v2 template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/v2 .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

## Структура проекта

Приведённая в этом репозитории структура проекта является рекомендуемой, но не обязательной.

Это лишь пример организации кода, который поможет вам в реализации сервиса.

При необходимости можно вносить изменения в структуру проекта, использовать любые библиотеки и предпочитаемые структурные паттерны организации кода приложения, например:
- **DDD** (Domain-Driven Design)
- **Clean Architecture**
- **Hexagonal Architecture**
- **Layered Architecture**

## Вывод go tool pprof -top -sample_index=alloc_space -diff_base="base.pprof" "result.pprof" 

File: main.exe
Build ID: C:\Users\E784~1\AppData\Local\Temp\go-build2137909896\b001\exe\main.exe2026-06-17 13:54:13.392393 +0500 +05
Type: alloc_space
Time: 2026-06-17 12:52:25 +05
Duration: 60.02s, Total samples = 15207.36MB 
Showing nodes accounting for -2196.05MB, 14.44% of 15207.36MB total
Dropped 291 nodes (cum <= 76.04MB)
      flat  flat%   sum%        cum   cum%
-2567.18MB 16.88% 16.88% -2567.18MB 16.88%  github.com/dimalewshin98-glitch/ShortyURL/internal/repository.(*InmemoryRepository).GetUsersID
  275.01MB  1.81% 15.07%   342.85MB  2.25%  compress/flate.NewWriter (inline)
   70.84MB  0.47% 14.61%    70.84MB  0.47%  compress/flate.(*compressor).initDeflate (inline)
   24.24MB  0.16% 14.45%    24.24MB  0.16%  net/http.init.func16
   13.05MB 0.086% 14.36%    13.05MB 0.086%  bufio.NewReaderSize (inline)
  -12.51MB 0.082% 14.44%   -12.51MB 0.082%  sync.(*Pool).pinSlow
    2.50MB 0.016% 14.43% -2215.61MB 14.57%  main.mainn.RequestLogger.func3
      -2MB 0.013% 14.44%    15.53MB   0.1%  net/http.(*Transport).dialConn
      -1MB 0.0066% 14.45% -2233.75MB 14.69%  net/http.(*conn).serve
    0.50MB 0.0033% 14.44%   348.08MB  2.29%  github.com/dimalewshin98-glitch/ShortyURL/internal/handler.(*RequestsHandler).ApiShortenBatch
    0.50MB 0.0033% 14.44% -2215.62MB 14.57%  main.mainn.AuthMiddleware.func2
         0     0% 14.44%    67.84MB  0.45%  compress/flate.(*compressor).init
         0     0% 14.44%   342.85MB  2.25%  compress/gzip.(*Writer).Write
         0     0% 14.44%   341.97MB  2.25%  encoding/json.(*Encoder).Encode
         0     0% 14.44%   341.97MB  2.25%  github.com/dimalewshin98-glitch/ShortyURL/internal/handler.(*compressWriter).Write
         0     0% 14.44% -2567.18MB 16.88%  github.com/dimalewshin98-glitch/ShortyURL/internal/handler.CreateUserID
         0     0% 14.44%   352.08MB  2.32%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 14.44%   348.08MB  2.29%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0% 14.44%    11.50MB 0.076%  main.loadtestfunc
         0     0% 14.44%   343.07MB  2.26%  main.mainn.GzipMiddleware.func1
         0     0% 14.44%       11MB 0.072%  net.(*netFD).dial
         0     0% 14.44%       12MB 0.079%  net.(*sysDialer).dialParallel.func1
         0     0% 14.44%       12MB 0.079%  net.(*sysDialer).dialSerial
         0     0% 14.44%       11MB 0.072%  net.(*sysDialer).dialTCP
         0     0% 14.44%       11MB 0.072%  net.(*sysDialer).doDialTCP (inline)
         0     0% 14.44%       11MB 0.072%  net.(*sysDialer).doDialTCPProto
         0     0% 14.44%       11MB 0.072%  net.internetSocket
         0     0% 14.44%       11MB 0.072%  net.socket
         0     0% 14.44%    23.24MB  0.15%  net/http.(*Request).write
         0     0% 14.44%    15.53MB   0.1%  net/http.(*Transport).dialConnFor
         0     0% 14.44%    15.53MB   0.1%  net/http.(*Transport).startDialConnForLocked.func1
         0     0% 14.44%    24.74MB  0.16%  net/http.(*persistConn).writeLoop
         0     0% 14.44%    22.74MB  0.15%  net/http.(*transferWriter).doBodyCopy
         0     0% 14.44%    21.74MB  0.14%  net/http.(*transferWriter).writeBody
         0     0% 14.44% -2214.74MB 14.56%  net/http.HandlerFunc.ServeHTTP
         0     0% 14.44%    22.74MB  0.15%  net/http.getCopyBuf (inline)
         0     0% 14.44% -2214.74MB 14.56%  net/http.serverHandler.ServeHTTP
         0     0% 14.44%    26.74MB  0.18%  sync.(*Pool).Get
         0     0% 14.44%   -11.51MB 0.076%  sync.(*Pool).Put
         0     0% 14.44%   -12.51MB 0.082%  sync.(*Pool).pin