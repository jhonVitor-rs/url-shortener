// Package infra providencia implementações de infraestrutura como cache
package infra

import (
	"context"
	"log/slog"
	"time"

	"github.com/jhonVitor-rs/url-shortener/internal/core/domain/models"
	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
	"github.com/redis/go-redis/v9"
)

const (
	listKey    = "url:recent"    // Chave para lista de URLs recentes
	maxLength  = 20              // Tamanho máximo da lista
	urlPrefix  = "url:"          // Prefixo para as chaves de URL no Redis
	defaultTTL = 24 * time.Hour  // TTL padrão para URLs sem data de expiração
	minTTL     = 5 * time.Minute // TTL mínimo para evitar expiração imediata
)

// URLCache encapsula funcionalidades de cache para URLs
type URLCache struct {
	client *redis.Client
	logger *slog.Logger
}

// NewURLCache cria uma nova instância do cache de URLs
func NewURLCache(client *redis.Client) *URLCache {
	return &URLCache{
		client: client,
		logger: slog.Default().With("component", "url_cache"),
	}
}

// LogRecentAccess armazena a URL original em cache e atualiza a lista de acessos recentes
// Esta função executa assincronamente e não bloqueia o chamador
func (c *URLCache) LogRecentAccess(ctx context.Context, shortUrl *models.ShortUrl) error {
	if shortUrl == nil {
		return wraperrors.ValidationErr("short url cannot be nil")
	}

	go func() {
		c.logger.Info("log recent access started", "slug", shortUrl.Slug)

		// Usar contexto Background para operação em goroutine separada
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		urlKey := urlPrefix + shortUrl.Slug

		// Determinar TTL com base na data de expiração
		ttl := c.calculateTTL(shortUrl)
		if ttl <= 0 {
			c.logger.Debug("url already expired, not caching", "slug", shortUrl.Slug)
			return
		}

		// Armazenar URL em cache
		if _, err := c.client.Set(timeoutCtx, urlKey, shortUrl.OriginalUrl, ttl).Result(); err != nil {
			c.logger.Error("failed to save url in cache", "error", err, "slug", shortUrl.Slug)
			return
		}

		c.logger.Info("url cached successfully", "slug", shortUrl.Slug, "ttl", ttl)

		// Atualizar lista de URLs recentes
		c.updateListURLs(timeoutCtx, shortUrl.Slug)
	}()

	return nil
}

// GetUrl recupera uma URL do cache pelo slug
func (c *URLCache) GetURL(ctx context.Context, slug string) (string, error) {
	if slug == "" {
		return "", wraperrors.ValidationErr("slug cannot be empty")
	}

	urlKey := urlPrefix + slug
	url, err := c.client.Get(ctx, urlKey).Result()
	if err != nil {
		// Melhorar diferenciação entre erros
		if err == redis.Nil {
			return "", wraperrors.NotFoundErr("url not found in cache")
		}

		c.logger.Error("failed to retrieve url from cache", "slug", slug, "error", err)
		return "", wraperrors.InternalErr("cache retrieval error", err)
	}

	// Atualizar lista de URLs recentes de forma assíncrona
	go func() {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		c.updateListURLs(timeoutCtx, slug)
	}()

	return url, nil
}

// updateListURLs atualiza a lista de URLs recentes
func (c *URLCache) updateListURLs(ctx context.Context, slug string) {
	c.logger.Debug("updating cache list", "slug", slug)

	pipe := c.client.TxPipeline()
	pipe.LRem(ctx, listKey, 0, slug) // Remover slug se já existir
	pipe.LPush(ctx, listKey, slug)   // Adicionar slug no topo da lista
	pipe.LLen(ctx, listKey)          // Obter tamanho atual da lista

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		c.logger.Error("failed to update recent url list", "error", err, "slug", slug)
		return
	}

	// Verificar se precisamos limitar o tamanho da lista
	listLen := cmds[2].(*redis.IntCmd).Val()
	if listLen > maxLength {
		cleanPipe := c.client.Pipeline()
		cleanPipe.LTrim(ctx, listKey, 0, maxLength-1)
		if _, err := cleanPipe.Exec(ctx); err != nil {
			c.logger.Error("failed to trim recent url list", "error", err)
		}
	}
}

// calculateTTL calcula o TTL apropriado para o cache
func (c *URLCache) calculateTTL(shortUrl *models.ShortUrl) time.Duration {
	if shortUrl.ExpiresAt != nil && !shortUrl.ExpiresAt.IsZero() {
		ttl := time.Until(*shortUrl.ExpiresAt)

		// Definir um limite mínimo para o TTL para evitar expirações muito rápidas
		if ttl <= 0 {
			return 0
		}
		if ttl < minTTL {
			return minTTL
		}
		return ttl
	}

	return defaultTTL
}

// GetRecentURLs retorna a lista de slugs de URLs recentes
func (c *URLCache) GetRecentURLs(ctx context.Context, limit int) ([]string, error) {
	if limit <= 0 || limit > maxLength {
		limit = maxLength
	}

	slugs, err := c.client.LRange(ctx, listKey, 0, int64(limit-1)).Result()
	if err != nil {
		c.logger.Error("failed to get recent urls", "error", err)
		return nil, wraperrors.InternalErr("failed to retrieve recent urls", err)
	}

	return slugs, nil
}
