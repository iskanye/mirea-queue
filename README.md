<p align="center">
    <a href="https://github.com/iskanye/mirea-queue">
        <img width="200px" height="200px" alt="mirea-queue" src="docs/bot.webp">
    </a>
</p>

<h1 align="center">
    Бот очереди для РТУ МИРЭА
</h1>

<div align="center">
    
[![License](https://img.shields.io/github/license/iskanye/mirea-queue)](https://github.com/iskanye/mirea-queue/blob/main/LICENSE)
[![Deploy Status](https://img.shields.io/github/actions/workflow/status/iskanye/mirea-queue/deploy.yml)](https://github.com/iskanye/mirea-queue/actions)
    
</div>

<p align="center">
    <b>
        Для тех кому надоели бесконечные очереди, которые противоречат друг другу
    </b>
</p>

<h2>Функционал бота</h2>
<ul>
    <li>Позволяет записываться на очередь по предмету</li>
    <li>Разделение очередей по группам и предметам каждой группы</li>
    <li>Возможность пропускать человека выше по очереди</li>
    <li>Интеграция с расписанием МИРЭА (подгружает оттуда список групп и их дисциплин)</li>
    <li>Очищает все очереди в 6 утра по МСК (перед началом всех пар)</li>
</ul>

<h2>Установка и запуск</h2>

Для установки бота скопируйте исходный код бота, измените токен бота на токен того телеграмм бота, на котором будет работать очередь, и измените токен админа на любую строчку по желанию в файле [.env.example](/.env.example) и переименуйте его в `.env`

Для работы бота необходим установленный [Docker](https://www.docker.com/). После настройки запустите следующие команды:

```bash
docker compose -f docker-compose.local.yaml up -d --build
```

Пользоваться ботом можно по [ссылке](https://t.me/rtumirea_queuebot)
