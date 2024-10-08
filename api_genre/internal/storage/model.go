package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

type Genre struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	TranslitName string             `bson:"translit"`
	BookCount    int64              `bson:"bookCount"`
}
type Category struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	TranslitName string             `bson:"translit"`
	Genres       []*Genre           `bson:"genres"`
}

var mocked_categories = []struct {
	Name   string
	Genres []string
}{
	{
		Name: "Бизнес-книги",
		Genres: []string{
			"Менеджмент",
			"Работа с клиентами",
			"Стартапы и создание бизнеса",
			"Переговоры",
			"Ораторское искусство / риторика",
			"Тайм-менеджмент",
			"Личная эффективность",
			"Продажи",
			"Интернет-бизнес",
			"Зарубежная деловая литература",
			"Делопроизводство",
			"Малый и средний бизнес",
			"О бизнесе популярно",
			"Недвижимость",
			"Личные финансы",
			"Корпоративная культура",
			"Отраслевые издания",
			"Финансы",
			"Экономика",
			"Бухучет / налогообложение / аудит",
			"Ценные бумаги / инвестиции",
			"Банковское дело",
			"Маркетинг, PR, реклама",
			"Логистика",
			"Кадровый менеджмент",
			"Поиск работы / карьера",
			"Менеджмент и кадры",
			"Государственное и муниципальное управление",
			"Политическое управление",
			"Краткое содержание",
			"Бизнес-справочники",
		},
	},
	{
		Name: "Знания и навыки",
		Genres: []string{
			"Научно-популярная литература",
			"Учебная и научная литература",
			"Компьютерная литература",
			"Культура и искусство",
			"Саморазвитие / личностный рост",
			"Эзотерика",
			"Словари, справочники",
			"Путеводители",
			"Истории из жизни",
			"Изучение языков",
		},
	},
	{
		Name: "Хобби, досуг",
		Genres: []string{
			"Отдых / туризм",
			"Хобби / увлечения",
			"Охота",
			"Мода и стиль",
			"Автомобили и ПДД",
			"Сад и огород",
			"Прикладная литература",
			"Развлечения",
			"Рукоделие и ремесла",
			"Искусство фотографии",
			"Спорт / фитнес",
			"Изобразительное искусство",
			"Сделай сам",
			"Йога",
			"Кулинария",
			"Путеводители",
			"Природа и животные",
			"Рыбалка",
			"Интеллектуальные игры",
		},
	},
	{
		Name: "Легкое чтение",
		Genres: []string{
			"Детективы",
			"Фантастика",
			"Фэнтези",
			"Любовные романы",
			"Ужасы / мистика",
			"Боевики, остросюжетная литература",
			"Юмористическая литература",
			"Приключения",
			"Young adult",
			"Классика жанра",
			"Легкая проза",
		},
	},
	{
		Name: "История",
		Genres: []string{
			"Историческое фэнтези",
			"Исторические приключения",
			"Книги о войне",
			"Книги о путешествиях",
			"Исторические любовные романы",
			"Документальная литература",
			"Историческая литература",
			"Биографии и мемуары",
			"Историческая фантастика",
			"Морские приключения",
			"Исторические детективы",
			"Популярно об истории",
		},
	},
	{
		Name: "Дом, дача",
		Genres: []string{
			"Отдых / туризм",
			"Интерьеры",
			"Хобби / увлечения",
			"Охота",
			"Фэншуй / фэн-шуй",
			"Автомобили и ПДД",
			"Сад и огород",
			"Рукоделие и ремесла",
			"Домашние животные",
			"Сделай сам",
			"Кулинария",
			"Природа и животные",
			"Ремонт в квартире",
			"Домашнее хозяйство",
			"Рыбалка",
			"Комнатные растения",
		},
	},
	{
		Name: "Детские книги",
		Genres: []string{
			"Зарубежные детские книги",
			"Детские стихи",
			"Детские детективы",
			"Детская фантастика",
			"Детские приключения",
			"Сказки",
			"Школьные учебники",
			"Книги для подростков",
			"Буквари",
			"Детская проза",
			"Учебная литература",
			"Внеклассное чтение",
			"Детская познавательная и развивающая литература",
			"Книги для детей",
			"Книги для дошкольников",
		},
	},
}
