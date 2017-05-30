package api

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// DefaultLimit is the number of results we will return per page if the user doesn't specify another amount
	DefaultLimit = 25
	// DefaultLimitString is DefaultLimit but in string form because types are a thing
	DefaultLimitString = "25"
	// MaxLimit is the maximum number of objects Dairycart will return in a response
	MaxLimit = 50
)

///////////////////////////////////////////////////////////////////////////////
//                                                                           //
//                                                 .---.                     //
//      ____________________                  _.--(\,/\,`\                   //
//     |                   |                ,' . . .' . \/`-._               //
//     |      Helper       |                : . . .' .  . |  _`--.           //
//     |                   |                :,----'--.  . ; (  \  `-.        //
//     |     Functions     |              ,;:|,-----.`,   `  `\_)  . `-.     //
//     |____     __________|            ,;;;`;        /; ,         '   `.    //
//          \   /                       ,';;,-'`.    ,-;;;;/-.         .  :  //
//           \  /                     ,;;;,'     `-,;;;;;-'  `-      ; '     //
//            \ /                    .;;,'       ,;;;;-'   :  .     -'       //
//             \/                   .;.'       ,';;;':     `  `              //
//              \                  ;;'      .;;;-'  ;  ,             ,       //
//                ____,-------.__ `'      .','    ;   `           , '        //
//             ,-'     __,---.__ `-.      `'   ,-'  ,             '          //
//          _,'___,-    `\_:_/'  - `:       ,-'     `  ,                     //
//        ,'  ((_); . .    :   . .  :    ,-'           `                     //
//      ,' :  '   :  . . _.;._  .  ,' ,-'                                    //
//     ,'  `      `.  _,:,---.`-._,' ,              _,---------._            //
//    '      :     `.::||     ::||  '            ,-'         ___ `-.         //
//      `    `-.     ::||     ::||`-.     __,---.,--_      ,'---`, `.        //
//        ,    `     ::|| `-' ::||   `_,-' `\_:_/'   `--._ `-(_)-'  |        //
//        `          ::||      :||   ; .  .   :    .  .   ;       `.|        //
//                   ::||`-'   :||  : .  .  __`.__  .  .  :          |       //
//                   `:||      :|;  `.  .  ; ____ :  .  . `,         |       //
//                    :;       `'     `.  /;'    `|---._   ;          |      //
//                    `               ; );::.,   , \:::||/'           `.     //
//                                   .' :::||` - '  :::||              :     //
//                                  ,'  :::||       :::||              :.    //
//       (         `.   ,'         ,'   :::|| `--'  :::||  -'          ::    //
//         -.       `---'       _,-'    `::|:       :::||              `|    //
//          `-  (             ,'   ,     ::|| `   ' :::|| -'            |    //
//        )    )  -.         ,'    `     ::||       `::||   --'         :    //
//       ) -.      `     __,-'    ,       :||        ::||               :    //
//          `   _____,--'         `        \;        `;/'            : .:    //
//     __,-----'            ,        ` .                              `::    //
//    -'           ,         `      :   `  .         -                   :   //
//         ,    , `,     ,         `.__    `     -      -      ;      '::    //
//    (  , '    '  `     `    ;    ,'  `---.__            ;    ',      :     //
//     \ '  ,          ,     ,'               `-._         ;    ';    ::     //
//      `-. '          `    :   ;                 `.        ,    ;,   ::     //
//         )        ,       :   `;                  `.   ,' ` ,  ,'   ;      //
//    `.  )__    ,  `  ,    `.   `; __                :  ,    ` ,'   :       //
//     )  )  `---`____ `      `   _'--`-.             ;___,----:    ,        //
//    __ \)))          `-----. _,-' --   `------.__,  ,'       |    ;        //
//                         ,-' -__,_______           ;        :       ,'     //
//                          `--'          `---------'          `-----'       //
//                                                                           //
///////////////////////////////////////////////////////////////////////////////

// borrowed from http://stackoverflow.com/questions/32825640/custom-marshaltext-for-golang-sql-null-types

// There's not really a great solution for these two stinkers here. Because []byte is what's expected, passing
// nil results in an empty string. The original has []byte("null"), which I think is actually worse. At least
// an empty string is falsy in most languages. ¯\_(ツ)_/¯

// NullFloat64 is a json.Marshal-able 64-bit float.
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalText satisfies the encoding.TestMarshaler interface
func (nf NullFloat64) MarshalText() ([]byte, error) {
	if nf.Valid {
		nfv := nf.Float64
		return []byte(strconv.FormatFloat(nfv, 'f', -1, 64)), nil
	}
	return nil, nil
}

// UnmarshalText is a function which unmarshals a NullFloat64
func (nf *NullFloat64) UnmarshalText(text []byte) (err error) {
	s := string(text)
	nf.NullFloat64.Float64, err = strconv.ParseFloat(s, 64)
	nf.NullFloat64.Valid = err == nil
	// returning nil because we've ensured that Float64 is set to at least zero.
	return nil
}

// This isn't borrowed, but rather inferred from stuff I borrowed above

// NullString is a json.Marshal-able String.
type NullString struct {
	sql.NullString
}

// MarshalText satisfies the encoding.TestMarshaler interface
func (ns NullString) MarshalText() ([]byte, error) {
	if ns.Valid {
		nsv := ns.String
		return []byte(nsv), nil
	}
	return nil, nil
}

// UnmarshalText is a function which unmarshals a NullString so that gorilla/schema can parse it
func (ns *NullString) UnmarshalText(text []byte) (err error) {
	ns.String = string(text)
	ns.Valid = true
	return nil
}

// Round borrowed from https://gist.github.com/DavidVaini/10308388
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////
//    ¸,ø¤º°º¤ø,¸¸,ø¤º°   Everything after this point is not borrowed.   °º¤ø,¸¸,ø¤º°º¤ø,¸    //
////////////////////////////////////////////////////////////////////////////////////////////////

// ListResponse is a generic list response struct containing values that represent
// pagination, meant to be embedded into other object response structs
type ListResponse struct {
	Count uint64 `json:"count"`
	Limit uint8  `json:"limit"`
	Page  uint64 `json:"page"`
}

// QueryFilter represents a query filter
type QueryFilter struct {
	Page          uint64
	Limit         uint8
	CreatedAfter  time.Time
	CreatedBefore time.Time
	UpdatedAfter  time.Time
	UpdatedBefore time.Time
}

func parseRawFilterParams(rawFilterParams url.Values) *QueryFilter {
	qf := &QueryFilter{
		Page:  1,
		Limit: 25,
	}

	page := rawFilterParams["page"]
	if len(page) == 1 {
		i, err := strconv.ParseUint(page[0], 10, 64)
		if err != nil {
			log.Printf("encountered error when trying to parse query filter param %s: %v", `Page`, err)
		} else {
			qf.Page = i
		}
	}

	limit := rawFilterParams["limit"]
	if len(limit) == 1 {
		i, err := strconv.ParseFloat(limit[0], 64)
		if err != nil {
			log.Printf("encountered error when trying to parse query filter param %s: %v", `Limit`, err)
		} else {
			qf.Limit = uint8(math.Min(i, MaxLimit))
		}
	}

	updatedAfter := rawFilterParams["updated_after"]
	if len(updatedAfter) == 1 {
		i, err := strconv.ParseUint(updatedAfter[0], 10, 64)
		if err != nil {
			log.Printf("encountered error when trying to parse query filter param %s: %v", `UpdatedAfter`, err)
		} else {
			qf.UpdatedAfter = time.Unix(int64(i), 0)
		}
	}

	updatedBefore := rawFilterParams["updated_before"]
	if len(updatedBefore) == 1 {
		i, err := strconv.ParseUint(updatedBefore[0], 10, 64)
		if err != nil {
			log.Printf("encountered error when trying to parse query filter param %s: %v", `UpdatedBefore`, err)
		} else {
			qf.UpdatedBefore = time.Unix(int64(i), 0)
		}
	}

	createdAfter := rawFilterParams["created_after"]
	if len(createdAfter) == 1 {
		i, err := strconv.ParseUint(createdAfter[0], 10, 64)
		if err != nil {
			log.Printf("encountered error when trying to parse query filter param %s: %v", `CreatedAfter`, err)
		} else {
			qf.CreatedAfter = time.Unix(int64(i), 0)
		}
	}

	createdBefore := rawFilterParams["created_before"]
	if len(createdBefore) == 1 {
		i, err := strconv.ParseUint(createdBefore[0], 10, 64)
		if err != nil {
			log.Printf("encountered error when trying to parse query filter param %s: %v", `CreatedBefore`, err)
		} else {
			qf.CreatedBefore = time.Unix(int64(i), 0)
		}
	}

	return qf
}

// rowExistsInDB will return whether or not a product/attribute/etc with a given identifier exists in the database
func rowExistsInDB(db *sql.DB, table, identifier, id string) (bool, error) {
	var exists string

	query := buildRowExistenceQuery(table, identifier, id)
	err := db.QueryRow(query, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists == "true", err
}

func respondThatRowDoesNotExist(req *http.Request, res http.ResponseWriter, itemType, id string) {
	itemTypeToIdentifierMap := map[string]string{
		"product attribute":       "id",
		"product attribute value": "id",
		"product progenitor":      "id",
		"product":                 "sku",
	}

	// in case we forget one, default to ID
	identifier := itemTypeToIdentifierMap[itemType]
	if _, ok := itemTypeToIdentifierMap[itemType]; !ok {
		identifier = "identified by"
	}

	log.Printf("informing user that the %s they were looking for (%s %s) does not exist", itemType, identifier, id)
	http.Error(res, fmt.Sprintf("The %s you were looking for (%s `%s`) does not exist", itemType, identifier, id), http.StatusNotFound)
}

func notifyOfInvalidRequestBody(res http.ResponseWriter, err error) {
	log.Printf("Encountered this error decoding a request body: %v", err)
	http.Error(res, err.Error(), http.StatusBadRequest)
}

func notifyOfInternalIssue(res http.ResponseWriter, err error, attemptedTask string) {
	log.Println(fmt.Sprintf("Encountered this error trying to %s: %v", attemptedTask, err))
	http.Error(res, "Unexpected internal error", http.StatusInternalServerError)
}
