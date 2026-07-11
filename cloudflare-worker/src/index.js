const ORIGIN = "https://origin-api.cauai.fun";
const CACHE_VERSION = "20260616-keys-theme-fix-v1";
const STATIC_PATHS = new Set(["/logo.png", "/favicon.ico"]);
const HTML_CACHE_TTL_SECONDS = 300;
const API_CACHE_TTL_SECONDS = 300;
const HASHED_ASSET_CACHE_TTL_SECONDS = 31536000;
const STATIC_CACHE_TTL_SECONDS = 86400;
const PUBLIC_SETTINGS_PATH = "/api/v1/settings/public";
const PUBLIC_CACHE_PATHS = new Set([PUBLIC_SETTINGS_PATH]);
const STATIC_CONTENT_TYPES = [
  "application/javascript",
  "application/json",
  "font/",
  "image/",
  "text/css",
  "text/javascript",
];

function hasAuthenticatedContext(request) {
  return Boolean(request.headers.get("authorization") || request.headers.get("cookie"));
}

function makeCacheKey(url) {
  const versionedUrl = new URL(url.toString());
  versionedUrl.searchParams.set("__kqs_cache_v", CACHE_VERSION);
  return new Request(versionedUrl.toString(), { method: "GET" });
}

function makeCacheableResponse(response, ttlSeconds) {
  const headers = new Headers(response.headers);
  headers.delete("set-cookie");
  headers.set("cache-control", `public, max-age=${ttlSeconds}`);
  headers.set("x-sub2api-edge-cache-version", CACHE_VERSION);
  return new Response(response.clone().body, {
    status: response.status,
    statusText: response.statusText,
    headers,
  });
}

function makeClientCacheResponse(response, ttlSeconds, immutable = false, edgeCache = "MISS") {
  const headers = new Headers(response.headers);
  headers.delete("set-cookie");
  headers.delete("content-length");
  headers.set("cache-control", `public, max-age=${ttlSeconds}${immutable ? ", immutable" : ""}`);
  headers.set("x-sub2api-edge-cache", edgeCache);
  headers.set("x-sub2api-edge-cache-version", CACHE_VERSION);
  return new Response(response.clone().body, {
    status: response.status,
    statusText: response.statusText,
    headers,
  });
}

function isStaticAssetResponse(response) {
  const contentType = (response.headers.get("content-type") || "").toLowerCase();
  return STATIC_CONTENT_TYPES.some((allowedType) => contentType.startsWith(allowedType));
}

function staticCacheTTL(pathname) {
  return pathname.startsWith("/assets/") ? HASHED_ASSET_CACHE_TTL_SECONDS : STATIC_CACHE_TTL_SECONDS;
}

function shouldBypassCache(request) {
  const cacheControl = request.headers.get("cache-control") || "";
  const pragma = request.headers.get("pragma") || "";
  return /\b(no-cache|no-store)\b/i.test(cacheControl) || /\bno-cache\b/i.test(pragma);
}

async function cacheResponse(request, response, ttlSeconds, ctx) {
  if (response.headers.has("set-cookie")) {
    return;
  }
  const cacheKey = makeCacheKey(new URL(request.url));
  const cacheable = makeCacheableResponse(response, ttlSeconds);
  if (ctx?.waitUntil) {
    ctx.waitUntil(caches.default.put(cacheKey, cacheable));
  } else {
    await caches.default.put(cacheKey, cacheable);
  }
}

export default {
  async fetch(request, _env, ctx) {
    const incomingUrl = new URL(request.url);
    const isStaticAsset =
      incomingUrl.pathname.startsWith("/assets/") || STATIC_PATHS.has(incomingUrl.pathname);
    const isPublicCacheRequest =
      request.method === "GET" && PUBLIC_CACHE_PATHS.has(incomingUrl.pathname) && !hasAuthenticatedContext(request);
    const isCacheableHtmlRequest =
      request.method === "GET" &&
      !isStaticAsset &&
      !incomingUrl.pathname.startsWith("/api/") &&
      !hasAuthenticatedContext(request);

    if (request.method === "GET" && isStaticAsset && !shouldBypassCache(request)) {
      const cacheKey = makeCacheKey(incomingUrl);
      const cached = await caches.default.match(cacheKey);
      if (cached) {
        const ttlSeconds = staticCacheTTL(incomingUrl.pathname);
        return makeClientCacheResponse(cached, ttlSeconds, incomingUrl.pathname.startsWith("/assets/"), "HIT");
      }
    }

    if (isPublicCacheRequest && !shouldBypassCache(request)) {
      const cacheKey = makeCacheKey(incomingUrl);
      const cached = await caches.default.match(cacheKey);
      if (cached) {
        return makeClientCacheResponse(cached, API_CACHE_TTL_SECONDS, false, "HIT");
      }
    }

    if (isCacheableHtmlRequest && !shouldBypassCache(request)) {
      const cacheKey = makeCacheKey(incomingUrl);
      const cached = await caches.default.match(cacheKey);
      if (cached) {
        return makeClientCacheResponse(cached, HTML_CACHE_TTL_SECONDS, false, "HIT");
      }
    }

    const upstreamUrl = new URL(incomingUrl.pathname + incomingUrl.search, ORIGIN);

    const headers = new Headers(request.headers);
    const visitorIP = request.headers.get("cf-connecting-ip") || "";

    // Normalize visitor IP headers at the edge. In same-zone Worker subrequests,
    // Cloudflare derives CF-Connecting-IP from this Worker-controlled x-real-ip.
    headers.delete("cf-connecting-ip");
    headers.delete("x-real-ip");
    headers.delete("x-forwarded-for");
    if (visitorIP) {
      headers.set("x-real-ip", visitorIP);
      headers.set("x-forwarded-for", visitorIP);
    }
    headers.set("x-forwarded-host", incomingUrl.host);
    headers.set("x-forwarded-proto", "https");
    headers.delete("cf-ipcountry");
    headers.delete("cf-ray");
    headers.delete("cf-visitor");

    const init = {
      method: request.method,
      headers,
      redirect: "manual",
    };

    if (request.method !== "GET" && request.method !== "HEAD") {
      init.body = request.body;
    }

    if (isStaticAsset) {
      const response = await fetch(new Request(upstreamUrl.toString(), {
        ...init,
        cf: {
          cacheEverything: true,
          cacheTtl: staticCacheTTL(incomingUrl.pathname),
        },
      }));
      if (request.method === "GET" && response.ok && isStaticAssetResponse(response)) {
        const ttlSeconds = staticCacheTTL(incomingUrl.pathname);
        const clientResponse = makeClientCacheResponse(response, ttlSeconds, incomingUrl.pathname.startsWith("/assets/"), "MISS");
        await cacheResponse(request, clientResponse, ttlSeconds, ctx);
        return clientResponse;
      }
      return response;
    }

    if (isPublicCacheRequest) {
      const response = await fetch(new Request(upstreamUrl.toString(), init));
      if (response.ok) {
        const clientResponse = makeClientCacheResponse(response, API_CACHE_TTL_SECONDS, false, "MISS");
        await cacheResponse(request, clientResponse, API_CACHE_TTL_SECONDS, ctx);
        return clientResponse;
      }
      return response;
    }

    const response = await fetch(new Request(upstreamUrl.toString(), init));

    const responseHeaders = new Headers(response.headers);
    const contentType = responseHeaders.get("content-type") || "";

    if (!isStaticAsset && request.method === "GET" && contentType.includes("text/html")) {
      responseHeaders.delete("content-length");
      responseHeaders.delete("set-cookie");
      responseHeaders.set("content-type", "text/html; charset=utf-8");
      responseHeaders.set("x-sub2api-edge-cache-version", CACHE_VERSION);
      if (response.ok && isCacheableHtmlRequest) {
        responseHeaders.set("cache-control", `public, max-age=${HTML_CACHE_TTL_SECONDS}, stale-while-revalidate=60`);
        responseHeaders.set("x-sub2api-edge-cache", "MISS");
      }
      const clientResponse = new Response(response.clone().body, {
        status: response.status,
        statusText: response.statusText,
        headers: responseHeaders,
      });
      if (response.ok && isCacheableHtmlRequest) {
        await cacheResponse(request, clientResponse, HTML_CACHE_TTL_SECONDS, ctx);
      }
      return clientResponse;
    }

    return response;
  },
};
