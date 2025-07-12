package repository

// Repositories aggregates all repository interfaces for DI

type Repositories struct {
	Member MemberRepository
	// Add more repositories here as needed, e.g.:
	// Product ProductRepository
}
