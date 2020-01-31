# skolkovo-test

# 1. Притяните к себе все файлы
git clone https://github.com/archi-chester/skolkovo-test.git

# 2. Перейдите в папку с проектом
cd skolkovo-test

# 3. В файле skolkovo-test.conf настройки в JSON формате по умолчанию он смотрит на docker elasticsearch:9200 и отправляет запросы # партнерам на url заглушку, указанную в PaternURL
nano skolkovo-test.conf
CTRL + O

# 4. Запуск контейнеров
sudo docker-compose -f docker-compose.yml up -d --build
