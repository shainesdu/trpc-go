---
title: 'codec: add error handling for codec registration and retrieval'
---
# Introduction

This document will walk you through the recent changes made to the codec registration and retrieval process. The purpose of these changes is to enhance error handling, ensuring that codec operations are more robust and informative.

We will cover:

1. Why error handling was added to codec registration and retrieval.
2. How the error handling is implemented in the registration process.
3. How the error handling is implemented in the retrieval process.
4. The role of backward compatibility in the updated functions.

# Error handling in codec registration

The primary reason for adding error handling to codec registration is to prevent duplicate registrations and provide clear feedback when such attempts occur. This is crucial for maintaining the integrity of the codec registry and avoiding conflicts.

<SwmSnippet path="codec/codec.go" line="59">

---

The <SwmToken path="/codec/codec.go" pos="63:2:2" line-data="// Register defines the logic of register a codec by name. It will be">`Register`</SwmToken> function now checks if a codec with the given name already exists in either the server or client codec maps. If it does, an error is returned indicating the codec is already registered.

```
var ErrCodecAlreadyRegistered = errors.New("codec already registered")

var ErrCodecNotFound = errors.New("codec not found")

// Register defines the logic of register a codec by name. It will be
// called by init function defined by third package. If there is no server codec,
// the second param serverCodec can be nil.
func Register(name string, serverCodec Codec, clientCodec Codec) error {
	lock.Lock()
	defer lock.Unlock()
```

---

</SwmSnippet>

<SwmSnippet path="codec/codec.go" line="70">

---

The implementation ensures that both server and client codecs are checked for existing entries before proceeding with the registration.

```
	if _, serverExists := serverCodecs[name]; serverExists {
		return fmt.Errorf("%w: server codec with name '%s'", ErrCodecAlreadyRegistered, name)
	}
	if _, clientExists := clientCodecs[name]; clientExists {
		return fmt.Errorf("%w: client codec with name '%s'", ErrCodecAlreadyRegistered, name)
	}
	
	serverCodecs[name] = serverCodec
	clientCodecs[name] = clientCodec
	return nil
}
```

---

</SwmSnippet>

# Error handling in codec retrieval

Error handling in codec retrieval was introduced to provide informative feedback when a requested codec is not found. This helps users quickly identify issues related to missing codecs.

<SwmSnippet path="codec/codec.go" line="82">

---

The <SwmToken path="/codec/codec.go" pos="82:2:2" line-data="// GetServer returns the server codec by name.">`GetServer`</SwmToken> function now returns an error if the requested server codec does not exist, enhancing the clarity of the retrieval process.

```
// GetServer returns the server codec by name.
func GetServer(name string) (Codec, error) {
	lock.RLock()
	c, exists := serverCodecs[name]
	lock.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("%w: server codec with name '%s'", ErrCodecNotFound, name)
	}
	return c, nil
}
```

---

</SwmSnippet>

<SwmSnippet path="codec/codec.go" line="94">

---

Similarly, the <SwmToken path="/codec/codec.go" pos="94:2:2" line-data="// GetClient returns the client codec by name.">`GetClient`</SwmToken> function returns an error when the requested client codec is not found, ensuring consistent error handling across retrieval operations.

```
// GetClient returns the client codec by name.
func GetClient(name string) (Codec, error) {
	lock.RLock()
	c, exists := clientCodecs[name]
	lock.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("%w: client codec with name '%s'", ErrCodecNotFound, name)
	}
	return c, nil
}
```

---

</SwmSnippet>

# Backward compatibility

To maintain backward compatibility, versions of the registration and retrieval functions that ignore errors have been provided. These functions allow existing codebases to continue functioning without modification, while still benefiting from the updated logic.

<SwmSnippet path="codec/codec.go" line="106">

---

The <SwmToken path="/codec/codec.go" pos="106:2:2" line-data="// RegisterCompatible is a backward compatible version of Register that ignores errors.">`RegisterCompatible`</SwmToken> function calls the <SwmToken path="/codec/codec.go" pos="106:16:16" line-data="// RegisterCompatible is a backward compatible version of Register that ignores errors.">`Register`</SwmToken> function but ignores any errors, allowing for seamless integration with older systems.

```
// RegisterCompatible is a backward compatible version of Register that ignores errors.
func RegisterCompatible(name string, serverCodec Codec, clientCodec Codec) {
	_ = Register(name, serverCodec, clientCodec)
}

// GetServerCompatible is a backward compatible version of GetServer that ignores errors.
func GetServerCompatible(name string) Codec {
	c, _ := GetServer(name)
	return c
}
```

---

</SwmSnippet>

<SwmSnippet path="codec/codec.go" line="117">

---

The <SwmToken path="/codec/codec.go" pos="111:2:2" line-data="// GetServerCompatible is a backward compatible version of GetServer that ignores errors.">`GetServerCompatible`</SwmToken> and <SwmToken path="/codec/codec.go" pos="117:2:2" line-data="// GetClientCompatible is a backward compatible version of GetClient that ignores errors.">`GetClientCompatible`</SwmToken> functions similarly ignore errors, ensuring that existing retrieval operations remain unaffected by the new error handling.

```
// GetClientCompatible is a backward compatible version of GetClient that ignores errors.
func GetClientCompatible(name string) Codec {
	c, _ := GetClient(name)
```

---

</SwmSnippet>

These changes collectively enhance the robustness of codec operations while providing options for backward compatibility.

<SwmMeta version="3.0.0" repo-id="Z2l0aHViJTNBJTNBdHJwYy1nbyUzQSUzQXNoYWluZXNkdQ==" repo-name="trpc-go"><sup>Powered by [Swimm](https://app.swimm.io/)</sup></SwmMeta>
