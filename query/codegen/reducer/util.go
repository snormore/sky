package reducer

// submap retrieves the child of a map. If the child doesn't exist it is created.
func submap(m map[string]interface{}, key string) map[string]interface{} {
	if tmp, ok := m[key].(map[string]interface{}); ok {
		return tmp
	}

	tmp := make(map[string]interface{})
	m[key] = tmp
	return tmp
}
