### О проекте
mathbattle - это платформа для проведения соревнований по [математическим боям](https://ru.wikipedia.org/wiki/%D0%9C%D0%B0%D1%82%D0%B1%D0%BE%D0%B9). Включает в себя API server (mb-server) и Телеграмм бота (mb-bot)

### Сборка и использование

```bash
git clone https://github.com/tmilovanov/mathbattle.git
make
```

Скопируйте шаблон конфигурации в директорию со скомпилированным mbserver и mbbot

    cp config/config_template.yaml bin/config.yaml
    cd bin

Отредактируйте конфигурацию. В поле `token` укажите [telegram token](https://core.telegram.org/bots/api#authorizing-your-bot)

Запустите mbbot и mbserver

```
./mbserver &
./mbbot &
```