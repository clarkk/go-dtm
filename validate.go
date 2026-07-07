package dtm

import (
	"net"
	"context"
	"regexp"
	"strings"
	"unicode"
	"github.com/clarkk/go-api/errin"
	"github.com/clarkk/go-dbd"
	"github.com/clarkk/go-dtm/types"
	"github.com/clarkk/go-fmt/sanitize"
	"github.com/clarkk/go-util/hash_pass"
	"github.com/clarkk/go-util/rdb"
	"github.com/clarkk/go-util/secure_pass"
)

var (
	re_alphanum 		= regexp.MustCompile(`(?i)^[a-z\d]+$`)
	re_alphanum_chars 	= regexp.MustCompile(`(?i)^[a-z\d_.-]+$`)
	
	re_email_user 		= regexp.MustCompile(`(?i)^[a-z\d+_.-]+$`)
	re_email_domain 	= regexp.MustCompile(`^(?:[a-z\d](?:[a-z\d-]*[a-z\d])?\.)+[a-z]{2,}$`)
)

type (
	Validate struct {
		*model_base
	}
	
	exists_string_func		func(value string, omit_id uint64) (uint64, error)
	exists_uint64_func		func(value, omit_id uint64) (uint64, error)
)

func NewValidate(base *model_base) Validate {
	return Validate{
		model_base: base,
	}
}

func (m *Validate) Normalize(s string, allow_newlines bool) string {
	return sanitize.Trim(sanitize.Filter_utf8mb3(s), allow_newlines)
}

func (m *Validate) Normalize_ptr(s *string, allow_newlines bool){
	*s = m.Normalize(*s, allow_newlines)
}

func (m *Validate) Field_code(s *string){
	m.Normalize_ptr(s, false)
	*s = strings.ToUpper(*s)
}

func (m *Validate) Field_code_letters(field string, s *string) *errin.Lang {
	m.Field_code(s)
	for _, r := range *s {
		if r < 'A' || r > 'Z' {
			return &errin.Lang{"FIELD_ALPHA", errin.Rep{
				"field":	m.Env_lang_field(field),
			}}
		}
	}
	return nil
}

func (m *Validate) Field_table_maxlength(table, field string, s *string, allow_newlines bool){
	m.Normalize_ptr(s, allow_newlines)
	length := dbd.Schema(table, field).Length()
	v := []rune(*s)
	if length >= len(v) {
		return
	}
	*s = string(v[:length])
}

func (m *Validate) Error_empty(field string, s *string) *errin.Lang {
	if *s == "" {
		return &errin.Lang{"FIELD_EMPTY", errin.Rep{
			"field":	m.Env_lang_field(field),
		}}
	}
	return nil
}

func (m *Validate) Error_maxlength(field string, s *string) *errin.Lang {
	return m.Error_table_maxlength(m.Resource(), field, s)
}

func (m *Validate) Error_table_maxlength(table, field string, s *string) *errin.Lang {
	length := dbd.Schema(table, field).Length()
	if len([]rune(*s)) > length {
		return &errin.Lang{"FIELD_MAXLENGTH", errin.Rep{
			"field":	m.Env_lang_field(field),
			"length":	length,
		}}
	}
	return nil
}

func (m *Validate) Error_code(field string, s *string) *errin.Lang {
	m.Field_code(s)
	if err := m.Error_empty(field, s); err != nil {
		return err
	}
	if err := m.Error_maxlength(field, s); err != nil {
		return err
	}
	if !re_alphanum.MatchString(*s) {
		return &errin.Lang{"FIELD_ALPHA_NUM", errin.Rep{
			"field":	m.Env_lang_field(field),
		}}
	}
	return nil
}

func (m *Validate) Error_code_letters(field string, s *string) *errin.Lang {
	var err *errin.Lang
	if err = m.Field_code_letters(field, s); err != nil {
		return err
	}
	if err = m.Error_empty(field, s); err != nil {
		return err
	}
	if err = m.Error_maxlength(field, s); err != nil {
		return err
	}
	return nil
}

func (m *Validate) Error_code_chars(field string, s *string) *errin.Lang {
	m.Field_code(s)
	if err := m.Error_empty(field, s); err != nil {
		return err
	}
	if err := m.Error_maxlength(field, s); err != nil {
		return err
	}
	if !re_alphanum_chars.MatchString(*s) {
		return &errin.Lang{"FIELD_ALPHA_NUM_CHARS", errin.Rep{
			"field":	m.Env_lang_field(field),
			"chars":	strings.Join([]string{".","-","_"}, ""),
		}}
	}
	return nil
}

func (m *Validate) Error_string(field string, s *string) *errin.Lang {
	return m.Error_table_string(m.Resource(), field, s)
}

func (m *Validate) Error_table_string(table, field string, s *string) *errin.Lang {
	m.Normalize_ptr(s, false)
	if err := m.Error_empty(field, s); err != nil {
		return err
	}
	if err := m.Error_table_maxlength(table, field, s); err != nil {
		return err
	}
	return nil
}

func (m *Validate) Error_string_optional(field string, s *string) *errin.Lang {
	m.Normalize_ptr(s, false)
	if err := m.Error_maxlength(field, s); err != nil {
		return err
	}
	return nil
}

func (m *Validate) Error_digits(field string, s *string) *errin.Lang {
	m.Normalize_ptr(s, false)
	if err := m.Error_empty(field, s); err != nil {
		return err
	}
	if err := m.Error_maxlength(field, s); err != nil {
		return err
	}
	if !is_digits(*s) {
		return &errin.Lang{"FIELD_DIGITS", errin.Rep{
			"field":	m.Env_lang_field(field),
		}}
	}
	return nil
}

func (m *Validate) Error_table_string_optional(table, field string, s *string) *errin.Lang {
	m.Normalize_ptr(s, false)
	if err := m.Error_table_maxlength(table, field, s); err != nil {
		return err
	}
	return nil
}

func (m *Validate) Error_text_optional(field string, s *string) *errin.Lang {
	m.Normalize_ptr(s, true)
	if err := m.Error_maxlength(field, s); err != nil {
		return err
	}
	return nil
}

func (m *Validate) Error_string_duplicate(field string, s *string, omit_id uint64, fn exists_string_func) (string, *errin.Lang, error){
	exists, err := fn(*s, omit_id)
	if err != nil {
		return "", nil, err
	}
	if exists != 0 {
		return "", &errin.Lang{"FIELD_DUPLICATE", errin.Rep{
			"field":	m.Env_lang_field(field),
			"value":	*s,
		}}, nil
	}
	return *s, nil, nil
}

func (m *Validate) Error_uint64_duplicate(field string, i *uint64, omit_id uint64, fn exists_uint64_func) (uint64, *errin.Lang, error){
	exists, err := fn(*i, omit_id)
	if err != nil {
		return 0, nil, err
	}
	if exists != 0 {
		return 0, &errin.Lang{"FIELD_DUPLICATE", errin.Rep{
			"field":	m.Env_lang_field(field),
			"value":	*i,
		}}, nil
	}
	return *i, nil, nil
}

func (m *Validate) Error_string_uint64_duplicate(field string, i uint64, value string, omit_id uint64, fn exists_uint64_func) (uint64, *errin.Lang, error){
	exists, err := fn(i, omit_id)
	if err != nil {
		return 0, nil, err
	}
	if exists != 0 {
		return 0, &errin.Lang{"FIELD_DUPLICATE", errin.Rep{
			"field":	m.Env_lang_field(field),
			"value":	value,
		}}, nil
	}
	return i, nil, nil
}

func (m *Validate) Error_enum(field string, s *string, options []string) (int, *errin.Lang){
	value := *s
	if value != "" {
		for i, v := range options {
			if value == v {
				return i, nil
			}
		}
	}
	//	Remove first value if empty string
	if options[0] == "" {
		options = options[1:]
	}
	return 0, &errin.Lang{"FIELD_ENUM", errin.Rep{
		"field":	m.Env_lang_field(field),
		"options":	strings.Join(options, ", "),
	}}
}

func (m *Validate) Error_int_range(field string, i *int, min, max int64) *errin.Lang {
	i64 := int64(*i)
	if i64 < min || i64 > max {
		return &errin.Lang{"FIELD_NUMBER_RANGE", errin.Rep{
			"field":	m.Env_lang_field(field),
			"min":		min,
			"max":		max,
		}}
	}
	return nil
}

func (m *Validate) Error_uint64_range(field string, i *uint64, min, max int64) *errin.Lang {
	i64 := int64(*i)
	if i64 < min || i64 > max {
		return &errin.Lang{"FIELD_NUMBER_RANGE", errin.Rep{
			"field":	m.Env_lang_field(field),
			"min":		min,
			"max":		max,
		}}
	}
	return nil
}

func (m *Validate) Error_int64_range(field string, i, min, max, factor int64, allow_empty bool) *errin.Lang {
	if !allow_empty && i == 0 {
		return &errin.Lang{"FIELD_EMPTY", errin.Rep{
			"field":	m.Env_lang_field(field),
		}}
	}
	if i < min || i > max {
		return &errin.Lang{"FIELD_NUMBER_RANGE", errin.Rep{
			"field":	m.Env_lang_field(field),
			"min":		float64(min) / float64(factor),
			"max":		float64(max) / float64(factor),
		}}
	}
	return nil
}

func (m *Validate) Error_decimal(field string, s *dtp.Decimal, min, max, factor float64, allow_empty bool) (int64, *errin.Lang){
	i, err_fatal := s.Int64(factor)
	if err_fatal != nil {
		return 0, &errin.Lang{"FIELD_NUMBER_RANGE", errin.Rep{
			"field":	m.Env_lang_field(field),
			"min":		min / factor,
			"max":		max / factor,
		}}
	}
	if err := m.Error_int64_range(field, i, int64(min), int64(max), int64(factor), allow_empty); err != nil {
		return 0, err
	}
	return i, nil
}

func (m *Validate) Error_email(field string, s *string) (*errin.Lang, error){
	return m.Error_table_email(m.Resource(), field, s)
}

func (m *Validate) Error_table_email(table, field string, s *string) (*errin.Lang, error){
	if err := m.Error_table_string(table, field, s); err != nil {
		return err, nil
	}
	
	user, domain, found := strings.Cut(*s, "@")
	if !found {
		return &errin.Lang{"FIELD_EMAIL", errin.Rep{
			"email":	*s,
		}}, nil
	}
	domain = strings.ToLower(domain)
	if !re_email_user.MatchString(user) {
		return &errin.Lang{"FIELD_EMAIL_USER", errin.Rep{
			"user":		user,
			"email":	*s,
		}}, nil
	}
	if !re_email_domain.MatchString(domain) {
		return &errin.Lang{"FIELD_EMAIL_DOMAIN", errin.Rep{
			"domain":	domain,
			"email":	*s,
		}}, nil
	}
	key := "DNS-MX:"+domain
	ctx := context.Background()
	_, not_found, err := rdb.Get(ctx, key)
	if err != nil {
		if not_found {
			if mx, _ := net.LookupMX(domain); len(mx) == 0 {
				return &errin.Lang{"FIELD_EMAIL_DNSMX", errin.Rep{
					"email":	*s,
					"domain":	domain,
				}}, nil
			}
			rdb.Set(ctx, key, []byte("1"), 60*60*24*30)
		} else {
			return nil, err
		}
	}
	*s = user+"@"+domain
	return nil, nil
}

func (m *Validate) Error_password(field string, s *string) *errin.Lang {
	if err := m.Error_empty(field, s); err != nil {
		return err
	}
	if !secure_pass.Minimum_length(*s) {
		return &errin.Lang{"FIELD_PASSWORD_LENGTH", errin.Rep{
			"min":	secure_pass.MIN_LENGTH,
		}}
	}
	if !secure_pass.Entropy(*s) {
		return &errin.Lang{"FIELD_PASSWORD_ENTROPY", errin.Rep{
			"special":	secure_pass.SPECIAL_CHAR,
		}}
	}
	return nil
}

func (m *Validate) Error_password_confirm(field string, pass, pass_confirm *string) (string, *errin.Lang, error){
	if err := m.Error_empty(field, pass_confirm); err != nil {
		return "", err, nil
	}
	if *pass != *pass_confirm {
		return "", &errin.Lang{"FIELD_PASSWORD_CONFIRMATION", nil}, nil
	}
	s, err := hash_pass.Create(*pass_confirm)
	if err != nil {
		return "", nil, err
	}
	return s, nil, nil
}

func is_digits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}