package prekeyserver

import (
	"time"

	. "gopkg.in/check.v1"
)

func (s *GenericServerSuite) Test_clientProfile_shouldSerializeCorrectly(c *C) {
	cp := &clientProfile{}
	cp.identifier = 0xABCDEF11
	cp.instanceTag = 0x4253112A
	cp.publicKey = generateEDDSAPublicKeyFrom([symKeyLength]byte{0xAB, 0x42})
	cp.versions = []byte{0x04}
	cp.expiration = time.Date(2034, 11, 5, 13, 46, 00, 12, time.UTC)
	cp.dsaKey = nil
	cp.transitionalSignature = nil
	cp.sig = &eddsaSignature{
		s: [114]byte{0x15, 0x00, 0x00, 0x00, 0x12},
	}

	expected := []byte{
		0x0, 0x0, 0x0, 0x5,

		// identifier
		0x0, 0x1, 0xab, 0xcd, 0xef, 0x11,

		// instance tag
		0x0, 0x2, 0x42, 0x53, 0x11, 0x2a,

		// public key
		0x00, 0x03, 0x85, 0x9f, 0x37, 0x1f, 0xf3, 0x4f,
		0x36, 0x44, 0x5a, 0x99, 0xca, 0x8a, 0x11, 0x17,
		0x6b, 0xb8, 0x1e, 0xe0, 0x60, 0x39, 0x32, 0x76,
		0x71, 0xf4, 0xc6, 0x83, 0x77, 0x01, 0x45, 0x27,
		0x35, 0x3c, 0x75, 0xae, 0xee, 0xaa, 0xf9, 0x79,
		0x69, 0xa0, 0xd8, 0x9a, 0x3a, 0xb1, 0x48, 0xf6,
		0x44, 0x41, 0x83, 0x30, 0x9f, 0x41, 0x38, 0x1b,
		0xf3, 0x29, 0x00,

		// versions
		0x00, 0x05, 0x00, 0x00, 0x00, 0x01, 0x04,

		// expiry
		0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x79, 0xf8,
		0xc7, 0x98,

		// signature
		0x15, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00,
	}

	c.Assert(cp.serialize(), DeepEquals, expected)
}

func (s *GenericServerSuite) Test_serializeVersions_shouldSerializeCorrectly(c *C) {
	v := []byte{0x04, 0x03}

	expected := []byte{0x00, 0x00, 0x00, 0x02, 0x04, 0x03}

	c.Assert(serializeVersions(v), DeepEquals, expected)
}

func (s *GenericServerSuite) Test_serializeExpiry_shouldSerializeCorrectly(c *C) {
	t := time.Date(2034, 11, 5, 13, 46, 20, 12, time.UTC)

	expected := []byte{0x00, 0x00, 0x00, 0x00, 0x79, 0xf8, 0xc7, 0xac}

	c.Assert(serializeExpiry(t), DeepEquals, expected)
}

func (s *GenericServerSuite) Test_prekeyMessage_shouldSerializeCorrectly(c *C) {
	m := &prekeyMessage{}
	m.identifier = 0x4264212E
	m.instanceTag = 0x1234ABC0
	m.y = generateECDHPublicKeyFrom([symKeyLength]byte{0x42, 0x11, 0xAA, 0xDE, 0xAD, 0xBE, 0xEF})
	m.b = []byte{0x01, 0x03, 0x05, 0x07, 0x09, 0x0B, 0x0D, 0x0F}
	expected := []byte{
		// version
		0x00, 0x04,

		// message type
		0x0F,

		// identifier
		0x42, 0x64, 0x21, 0x2E,

		// instance tag
		0x12, 0x34, 0xAB, 0xC0,

		// y
		0x62, 0x38, 0x7d, 0xcd, 0x13, 0x84, 0x21, 0x0e,
		0x62, 0xcf, 0xaf, 0x06, 0x7f, 0x49, 0x02, 0x8c,
		0xdd, 0xfe, 0x99, 0xb9, 0x01, 0x59, 0x66, 0x7d,
		0x57, 0x0d, 0xc0, 0xb7, 0x89, 0x2c, 0xfc, 0x5c,
		0xac, 0xb8, 0x24, 0x17, 0xe9, 0x4d, 0x36, 0x29,
		0x04, 0x0e, 0x6a, 0xd1, 0xb4, 0x2d, 0x1a, 0x55,
		0xb9, 0x24, 0x29, 0x23, 0x7e, 0x5b, 0xc9, 0xe6,
		0x00,

		// b
		0x00, 0x00, 0x00, 0x08,
		0x01, 0x03, 0x05, 0x07, 0x09, 0x0B, 0x0D, 0x0F,
	}

	c.Assert(m.serialize(), DeepEquals, expected)
}

func (s *GenericServerSuite) Test_prekeyMessage_shouldDeserializeCorrectly(c *C) {
	m := &prekeyMessage{}
	_, ok := m.deserialize([]byte{
		// version
		0x00, 0x04,

		// message type
		0x0F,

		// identifier
		0x42, 0x64, 0x21, 0x2E,

		// instance tag
		0x12, 0x34, 0xAB, 0xC0,

		// y
		0x62, 0x38, 0x7d, 0xcd, 0x13, 0x84, 0x21, 0x0e,
		0x62, 0xcf, 0xaf, 0x06, 0x7f, 0x49, 0x02, 0x8c,
		0xdd, 0xfe, 0x99, 0xb9, 0x01, 0x59, 0x66, 0x7d,
		0x57, 0x0d, 0xc0, 0xb7, 0x89, 0x2c, 0xfc, 0x5c,
		0xac, 0xb8, 0x24, 0x17, 0xe9, 0x4d, 0x36, 0x29,
		0x04, 0x0e, 0x6a, 0xd1, 0xb4, 0x2d, 0x1a, 0x55,
		0xb9, 0x24, 0x29, 0x23, 0x7e, 0x5b, 0xc9, 0xe6,
		0x00,

		// b
		0x00, 0x00, 0x00, 0x08,
		0x01, 0x03, 0x05, 0x07, 0x09, 0x0B, 0x0D, 0x0F,
	})
	c.Assert(ok, Equals, true)
	c.Assert(m.identifier, Equals, uint32(0x4264212E))
	c.Assert(m.instanceTag, Equals, uint32(0x1234ABC0))
	c.Assert(m.y.k.Equals(generateECDHPublicKeyFrom([symKeyLength]byte{0x42, 0x11, 0xAA, 0xDE, 0xAD, 0xBE, 0xEF}).k), Equals, true)
	c.Assert(m.b, DeepEquals, []byte{0x01, 0x03, 0x05, 0x07, 0x09, 0x0B, 0x0D, 0x0F})
}

func (s *GenericServerSuite) Test_prekeyProfile_shouldSerializeCorrectly(c *C) {
	m := &prekeyProfile{}
	m.identifier = 0x4264212F
	m.instanceTag = 0x1234ABC1
	m.expiration = time.Date(2034, 11, 5, 13, 46, 00, 12, time.UTC)
	m.sharedPrekey = generateECDHPublicKeyFrom([symKeyLength]byte{0x44, 0x11, 0xAA, 0xDE, 0xAD, 0xBE, 0xEF})
	m.sig = &eddsaSignature{
		s: [114]byte{0x16, 0x00, 0x00, 0x00, 0x12, 0x11},
	}
	expected := []byte{
		// // version
		// 0x00, 0x04,

		// // message type
		// 0x0F,

		// identifier
		0x42, 0x64, 0x21, 0x2F,

		// instance tag
		0x12, 0x34, 0xAB, 0xC1,

		// expiration
		0x00, 0x00, 0x00, 0x00, 0x79, 0xf8, 0xc7, 0x98,

		// shared prekey
		0x82, 0xd7, 0xaf, 0x02, 0xa2, 0x05, 0xb6, 0x06,
		0x15, 0x2b, 0x9a, 0x83, 0x4e, 0x10, 0x33, 0xcc,
		0x64, 0x10, 0xaf, 0xce, 0x92, 0xa4, 0x35, 0x4f,
		0xc4, 0x67, 0x70, 0xc1, 0x5b, 0xec, 0x01, 0x5b,
		0xc4, 0x2e, 0xf9, 0x5a, 0x53, 0x06, 0x05, 0x50,
		0x51, 0x2a, 0x0a, 0xf2, 0xb5, 0x06, 0x4e, 0xac,
		0x88, 0x88, 0x69, 0x4f, 0xeb, 0x10, 0xef, 0x02,
		0x00,

		// sig
		0x16, 0x00, 0x00, 0x00, 0x12, 0x11, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
	}

	c.Assert(m.serialize(), DeepEquals, expected)
}

func (s *GenericServerSuite) Test_prekeyProfile_shouldDeserializeCorrectly(c *C) {
	m := &prekeyProfile{}
	_, ok := m.deserialize([]byte{
		// // version
		// 0x00, 0x04,

		// // message type
		// 0x0F,

		// identifier
		0x42, 0x64, 0x21, 0x2F,

		// instance tag
		0x12, 0x34, 0xAB, 0xC1,

		// expiration
		0x00, 0x00, 0x00, 0x00, 0x79, 0xf8, 0xc7, 0x98,

		// shared prekey
		0x82, 0xd7, 0xaf, 0x02, 0xa2, 0x05, 0xb6, 0x06,
		0x15, 0x2b, 0x9a, 0x83, 0x4e, 0x10, 0x33, 0xcc,
		0x64, 0x10, 0xaf, 0xce, 0x92, 0xa4, 0x35, 0x4f,
		0xc4, 0x67, 0x70, 0xc1, 0x5b, 0xec, 0x01, 0x5b,
		0xc4, 0x2e, 0xf9, 0x5a, 0x53, 0x06, 0x05, 0x50,
		0x51, 0x2a, 0x0a, 0xf2, 0xb5, 0x06, 0x4e, 0xac,
		0x88, 0x88, 0x69, 0x4f, 0xeb, 0x10, 0xef, 0x02,
		0x00,

		// sig
		0x16, 0x00, 0x00, 0x00, 0x12, 0x11, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,
	})
	c.Assert(ok, Equals, true)
	c.Assert(m.identifier, Equals, uint32(0x4264212F))
	c.Assert(m.instanceTag, Equals, uint32(0x1234ABC1))
	c.Assert(m.expiration, DeepEquals, time.Date(2034, 11, 5, 13, 46, 00, 00, time.UTC))
	c.Assert(m.sharedPrekey.k.Equals(generateECDHPublicKeyFrom([symKeyLength]byte{0x44, 0x11, 0xAA, 0xDE, 0xAD, 0xBE, 0xEF}).k), Equals, true)
	c.Assert(m.sig, DeepEquals, &eddsaSignature{
		s: [114]byte{0x16, 0x00, 0x00, 0x00, 0x12, 0x11},
	})
}

func (s *GenericServerSuite) Test_prekeyEnsemble_shouldSerializeCorrectly(c *C) {
	m := &prekeyEnsemble{}
	cp := &clientProfile{}
	cp.identifier = 0xABCDEF11
	cp.instanceTag = 0x4253112A
	cp.publicKey = generateEDDSAPublicKeyFrom([symKeyLength]byte{0xAB, 0x42})
	cp.versions = []byte{0x04}
	cp.expiration = time.Date(2034, 11, 5, 13, 46, 00, 12, time.UTC)
	cp.dsaKey = nil
	cp.transitionalSignature = nil
	cp.sig = &eddsaSignature{
		s: [114]byte{0x15, 0x00, 0x00, 0x00, 0x12},
	}
	m.cp = cp

	pp := &prekeyProfile{}
	pp.identifier = 0x4264212F
	pp.instanceTag = 0x1234ABC1
	pp.expiration = time.Date(2034, 11, 5, 13, 46, 00, 12, time.UTC)
	pp.sharedPrekey = generateECDHPublicKeyFrom([symKeyLength]byte{0x44, 0x11, 0xAA, 0xDE, 0xAD, 0xBE, 0xEF})
	pp.sig = &eddsaSignature{
		s: [114]byte{0x16, 0x00, 0x00, 0x00, 0x12, 0x11},
	}
	m.pp = pp

	pm := &prekeyMessage{}
	pm.identifier = 0x4264212E
	pm.instanceTag = 0x1234ABC0
	pm.y = generateECDHPublicKeyFrom([symKeyLength]byte{0x42, 0x11, 0xAA, 0xDE, 0xAD, 0xBE, 0xEF})
	pm.b = []byte{0x01, 0x03, 0x05, 0x07, 0x09, 0x0B, 0x0D, 0x0F}
	m.pm = pm

	expected := []byte{
		// - client profile

		0x0, 0x0, 0x0, 0x5,

		// identifier
		0x0, 0x1, 0xab, 0xcd, 0xef, 0x11,

		// instance tag
		0x0, 0x2, 0x42, 0x53, 0x11, 0x2a,

		// public key
		0x00, 0x03, 0x85, 0x9f, 0x37, 0x1f, 0xf3, 0x4f,
		0x36, 0x44, 0x5a, 0x99, 0xca, 0x8a, 0x11, 0x17,
		0x6b, 0xb8, 0x1e, 0xe0, 0x60, 0x39, 0x32, 0x76,
		0x71, 0xf4, 0xc6, 0x83, 0x77, 0x01, 0x45, 0x27,
		0x35, 0x3c, 0x75, 0xae, 0xee, 0xaa, 0xf9, 0x79,
		0x69, 0xa0, 0xd8, 0x9a, 0x3a, 0xb1, 0x48, 0xf6,
		0x44, 0x41, 0x83, 0x30, 0x9f, 0x41, 0x38, 0x1b,
		0xf3, 0x29, 0x00,

		// versions
		0x00, 0x05, 0x00, 0x00, 0x00, 0x01, 0x04,

		// expiry
		0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x79, 0xf8,
		0xc7, 0x98,

		// signature
		0x15, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00,

		// - prekey profile

		// // version
		// 0x00, 0x04,

		// // message type
		// 0x0F,

		// identifier
		0x42, 0x64, 0x21, 0x2F,

		// instance tag
		0x12, 0x34, 0xAB, 0xC1,

		// expiration
		0x00, 0x00, 0x00, 0x00, 0x79, 0xf8, 0xc7, 0x98,

		// shared prekey
		0x82, 0xd7, 0xaf, 0x02, 0xa2, 0x05, 0xb6, 0x06,
		0x15, 0x2b, 0x9a, 0x83, 0x4e, 0x10, 0x33, 0xcc,
		0x64, 0x10, 0xaf, 0xce, 0x92, 0xa4, 0x35, 0x4f,
		0xc4, 0x67, 0x70, 0xc1, 0x5b, 0xec, 0x01, 0x5b,
		0xc4, 0x2e, 0xf9, 0x5a, 0x53, 0x06, 0x05, 0x50,
		0x51, 0x2a, 0x0a, 0xf2, 0xb5, 0x06, 0x4e, 0xac,
		0x88, 0x88, 0x69, 0x4f, 0xeb, 0x10, 0xef, 0x02,
		0x00,

		// sig
		0x16, 0x00, 0x00, 0x00, 0x12, 0x11, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,

		// - prekey message

		// version
		0x00, 0x04,

		// message type
		0x0F,

		// identifier
		0x42, 0x64, 0x21, 0x2E,

		// instance tag
		0x12, 0x34, 0xAB, 0xC0,

		// y
		0x62, 0x38, 0x7d, 0xcd, 0x13, 0x84, 0x21, 0x0e,
		0x62, 0xcf, 0xaf, 0x06, 0x7f, 0x49, 0x02, 0x8c,
		0xdd, 0xfe, 0x99, 0xb9, 0x01, 0x59, 0x66, 0x7d,
		0x57, 0x0d, 0xc0, 0xb7, 0x89, 0x2c, 0xfc, 0x5c,
		0xac, 0xb8, 0x24, 0x17, 0xe9, 0x4d, 0x36, 0x29,
		0x04, 0x0e, 0x6a, 0xd1, 0xb4, 0x2d, 0x1a, 0x55,
		0xb9, 0x24, 0x29, 0x23, 0x7e, 0x5b, 0xc9, 0xe6,
		0x00,

		// b
		0x00, 0x00, 0x00, 0x08,
		0x01, 0x03, 0x05, 0x07, 0x09, 0x0B, 0x0D, 0x0F,
	}

	c.Assert(m.serialize(), DeepEquals, expected)
}

func (s *GenericServerSuite) Test_prekeyEnsemble_shouldDeserializeCorrectly(c *C) {
	m := &prekeyEnsemble{}
	_, ok := m.deserialize([]byte{
		// - client profile

		0x0, 0x0, 0x0, 0x5,

		// identifier
		0x0, 0x1, 0xab, 0xcd, 0xef, 0x11,

		// instance tag
		0x0, 0x2, 0x42, 0x53, 0x11, 0x2a,

		// public key
		0x00, 0x03, 0x85, 0x9f, 0x37, 0x1f, 0xf3, 0x4f,
		0x36, 0x44, 0x5a, 0x99, 0xca, 0x8a, 0x11, 0x17,
		0x6b, 0xb8, 0x1e, 0xe0, 0x60, 0x39, 0x32, 0x76,
		0x71, 0xf4, 0xc6, 0x83, 0x77, 0x01, 0x45, 0x27,
		0x35, 0x3c, 0x75, 0xae, 0xee, 0xaa, 0xf9, 0x79,
		0x69, 0xa0, 0xd8, 0x9a, 0x3a, 0xb1, 0x48, 0xf6,
		0x44, 0x41, 0x83, 0x30, 0x9f, 0x41, 0x38, 0x1b,
		0xf3, 0x29, 0x00,

		// versions
		0x00, 0x05, 0x00, 0x00, 0x00, 0x01, 0x04,

		// expiry
		0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x79, 0xf8,
		0xc7, 0x98,

		// signature
		0x15, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00,

		// - prekey profile

		// // version
		// 0x00, 0x04,

		// // message type
		// 0x0F,

		// identifier
		0x42, 0x64, 0x21, 0x2F,

		// instance tag
		0x12, 0x34, 0xAB, 0xC1,

		// expiration
		0x00, 0x00, 0x00, 0x00, 0x79, 0xf8, 0xc7, 0x98,

		// shared prekey
		0x82, 0xd7, 0xaf, 0x02, 0xa2, 0x05, 0xb6, 0x06,
		0x15, 0x2b, 0x9a, 0x83, 0x4e, 0x10, 0x33, 0xcc,
		0x64, 0x10, 0xaf, 0xce, 0x92, 0xa4, 0x35, 0x4f,
		0xc4, 0x67, 0x70, 0xc1, 0x5b, 0xec, 0x01, 0x5b,
		0xc4, 0x2e, 0xf9, 0x5a, 0x53, 0x06, 0x05, 0x50,
		0x51, 0x2a, 0x0a, 0xf2, 0xb5, 0x06, 0x4e, 0xac,
		0x88, 0x88, 0x69, 0x4f, 0xeb, 0x10, 0xef, 0x02,
		0x00,

		// sig
		0x16, 0x00, 0x00, 0x00, 0x12, 0x11, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00,

		// - prekey message

		// version
		0x00, 0x04,

		// message type
		0x0F,

		// identifier
		0x42, 0x64, 0x21, 0x2E,

		// instance tag
		0x12, 0x34, 0xAB, 0xC0,

		// y
		0x62, 0x38, 0x7d, 0xcd, 0x13, 0x84, 0x21, 0x0e,
		0x62, 0xcf, 0xaf, 0x06, 0x7f, 0x49, 0x02, 0x8c,
		0xdd, 0xfe, 0x99, 0xb9, 0x01, 0x59, 0x66, 0x7d,
		0x57, 0x0d, 0xc0, 0xb7, 0x89, 0x2c, 0xfc, 0x5c,
		0xac, 0xb8, 0x24, 0x17, 0xe9, 0x4d, 0x36, 0x29,
		0x04, 0x0e, 0x6a, 0xd1, 0xb4, 0x2d, 0x1a, 0x55,
		0xb9, 0x24, 0x29, 0x23, 0x7e, 0x5b, 0xc9, 0xe6,
		0x00,

		// b
		0x00, 0x00, 0x00, 0x08,
		0x01, 0x03, 0x05, 0x07, 0x09, 0x0B, 0x0D, 0x0F,
	})

	c.Assert(ok, Equals, true)

	c.Assert(m.cp.identifier, Equals, uint32(0xABCDEF11))
	c.Assert(m.cp.instanceTag, Equals, uint32(0x4253112A))
	c.Assert(m.cp.publicKey.k.Equals(generateEDDSAPublicKeyFrom([symKeyLength]byte{0xAB, 0x42}).k), Equals, true)
	c.Assert(m.cp.versions, DeepEquals, []byte{0x04})
	c.Assert(m.cp.expiration, Equals, time.Date(2034, 11, 5, 13, 46, 00, 00, time.UTC))
	c.Assert(m.cp.dsaKey, IsNil)
	c.Assert(m.cp.transitionalSignature, IsNil)
	c.Assert(m.cp.sig, DeepEquals, &eddsaSignature{
		s: [114]byte{0x15, 0x00, 0x00, 0x00, 0x12},
	})

	c.Assert(m.pp.identifier, Equals, uint32(0x4264212F))
	c.Assert(m.pp.instanceTag, Equals, uint32(0x1234ABC1))
	c.Assert(m.pp.expiration, DeepEquals, time.Date(2034, 11, 5, 13, 46, 00, 00, time.UTC))
	c.Assert(m.pp.sharedPrekey.k.Equals(generateECDHPublicKeyFrom([symKeyLength]byte{0x44, 0x11, 0xAA, 0xDE, 0xAD, 0xBE, 0xEF}).k), Equals, true)
	c.Assert(m.pp.sig, DeepEquals, &eddsaSignature{
		s: [114]byte{0x16, 0x00, 0x00, 0x00, 0x12, 0x11},
	})

	c.Assert(m.pm.identifier, Equals, uint32(0x4264212E))
	c.Assert(m.pm.instanceTag, Equals, uint32(0x1234ABC0))
	c.Assert(m.pm.y.k.Equals(generateECDHPublicKeyFrom([symKeyLength]byte{0x42, 0x11, 0xAA, 0xDE, 0xAD, 0xBE, 0xEF}).k), Equals, true)
	c.Assert(m.pm.b, DeepEquals, []byte{0x01, 0x03, 0x05, 0x07, 0x09, 0x0B, 0x0D, 0x0F})
}
