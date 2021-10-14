package leancloud

type UserRef struct {
	c     *Client
	class string
	ID    string
	token string
}

func (client *Client) User(user interface{}) *UserRef {
	if meta := extractUserMeta(user); meta != nil {
		return &UserRef{
			c:     client,
			class: "users",
			ID:    meta.ID,
		}
	}
	return nil
}

func (ref *Users) ID(id string) *UserRef {
	return &UserRef{
		c:     ref.c,
		class: "users",
		ID:    id,
	}
}

func (ref *UserRef) Get(user interface{}, authOptions ...AuthOption) error {
	if ref == nil || ref.ID == "" || ref.class == "" {
		return nil
	}
	if err := objectGet(ref, user, authOptions...); err != nil {
		return err
	}

	return nil
}

func (ref *UserRef) Set(key string, value interface{}, authOptions ...AuthOption) error {
	if ref == nil || ref.ID == "" || ref.class == "" {
		return nil
	}
	if err := objectSet(ref, key, value, authOptions...); err != nil {
		return err
	}

	return nil
}

func (ref *UserRef) Update(diff interface{}, authOptions ...AuthOption) error {
	if ref == nil || ref.ID == "" || ref.class == "" {
		return nil
	}
	if err := objectUpdate(ref, diff, authOptions...); err != nil {
		return err
	}

	return nil
}

func (ref *UserRef) UpdateWithQuery(diff interface{}, query *Query, authOptions ...AuthOption) error {
	if ref == nil || ref.ID == "" || ref.class == "" {
		return nil
	}
	if err := objectUpdateWithQuery(ref, diff, query, authOptions...); err != nil {
		return err
	}

	return nil
}

func (ref *UserRef) Destroy(authOptions ...AuthOption) error {
	if ref == nil || ref.ID == "" || ref.class == "" {
		return nil
	}
	if err := objectDestroy(ref, authOptions...); err != nil {
		return err
	}

	return nil
}
