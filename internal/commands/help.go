package commands

func Help(cn string) (string, error) {
	switch cn {
	case "phonetic":
		if resp, err := phonetic(); err != nil {
			return "", ErrInternalServerError
		} else {
			return resp, nil
		}
	}
	return "this shouldnt happen xD", nil
}

func phonetic() (string, error) {
	// This might look like complete ass depending on the
	// matrix clients font. Looks fine on my Element client.
	help := `
| Ё | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 0 | - | Ъ |
     | Я | Ш | Е | Р | Т | Ы | У | И | О | П | Ю | Щ | Э |
       | А | С | Д | Ф | Г | Ч | Й | К | Л | Ь | Ж |
           | З | Х | Ц | В | Б | Н | М | ; | : |
	`
	return help, nil
}
