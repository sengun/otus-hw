package hw03frequencyanalysis

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = true

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})
}

func TestPrepareText(t *testing.T) {
	dataSet := []struct {
		inputString    string
		expectedString string
	}{
		{inputString: "Привет, Отус- !", expectedString: "Привет  Отус !"},
		{inputString: "Пу-ух! Пу-ух!- а он  не  откликается,", expectedString: "Пу-ух! Пу-ух! а он  не  откликается "},
		{inputString: "Дефис-дефис -тире в начале или тире в конце- или - посередине", expectedString: "Дефис-дефис тире в начале или тире в конце или  посередине"},
		{inputString: "Запятая, точка.двоеточие: ", expectedString: "Запятая  точка двоеточие  "},
	}

	for i, ds := range dataSet {
		t.Run(fmt.Sprintf("Prepare Text with regular expressions. Dataset %d", i), func(t *testing.T) {
			require.Equal(t, ds.expectedString, PrepareText(ds.inputString))
		})
	}
}

func TestGetTextUnits(t *testing.T) {
	dataSet := []struct {
		inputString   string
		expectedSlice []string
	}{
		{inputString: "Привет, Отус- !", expectedSlice: []string{"Привет,", "Отус-", "!"}},
		{inputString: "", expectedSlice: []string{}},
	}

	for i, ds := range dataSet {
		t.Run(fmt.Sprintf("Split text. Dataset %d", i), func(t *testing.T) {
			require.Equal(t, ds.expectedSlice, GetTextUnits(ds.inputString))
		})
	}
}

func TestGetFirstTenWords(t *testing.T) {
	dataSet := []struct {
		inputStruct   []wordCountStruct
		expectedSlice []string
	}{
		{
			inputStruct:   []wordCountStruct{{"Привет,", 1}, {"Отус-", 1}, {"!", 1}},
			expectedSlice: []string{"Привет,", "Отус-", "!"},
		}, {
			inputStruct:   []wordCountStruct{{"1", 1}, {"2", 1}, {"3", 1}, {"4", 1}, {"5", 1}, {"6", 1}, {"7", 1}, {"8", 1}, {"9", 1}, {"10", 1}, {"11", 1}},
			expectedSlice: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
		}, {
			inputStruct:   []wordCountStruct{{"1", 1}, {"2", 1}, {"3", 1}, {"4", 1}, {"5", 1}, {"6", 1}, {"7", 1}, {"8", 1}, {"9", 1}, {"10", 1}},
			expectedSlice: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
		},
	}

	for i, ds := range dataSet {
		t.Run(fmt.Sprintf("Get First Ten Words. Dataset %d", i), func(t *testing.T) {
			require.Equal(t, ds.expectedSlice, GetFirstTenWords(ds.inputStruct))
		})
	}
}

func TestSortWordCountStruct(t *testing.T) {
	t.Run("Sort WordCount Struct", func(t *testing.T) {
		require.Equal(
			t,
			[]wordCountStruct{{"В", 3}, {"А", 1}, {"Б", 1}},
			SortWordCountStruct([]wordCountStruct{{"А", 1}, {"Б", 1}, {"В", 3}}),
		)
	})
}
