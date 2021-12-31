package jwt

var jwtFactory *JWTFactory

func SetJwtFactory(jf *JWTFactory) {
	jwtFactory = jf
}

func JwtFactory() *JWTFactory {
	return jwtFactory
}
