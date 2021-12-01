this.addEventListener("install", function(event) {
  console.log("Service Worker installing.");
  event.waitUntil(
      caches.open("V5")
          .then(function(cache) {
              cache.addAll(
                  ["/", "spectre.css", "elm.js", "manifest.json"]
              );
          })
  );
});
this.addEventListener("activate", function(event) {
    // delete any unexpected caches
    event.waitUntil(
        caches
            .keys()
            .then(function(keys) {
                return keys.filter(function(key) {
                    return key !== "V5";
                })
            })
            .then(function(keys) {
                Promise.all(
                    keys.map(function(key) {
                        console.log("Deleting cache " + key);
                        return caches.delete(key);
                    })
                )
            })
    );
});
this.addEventListener("fetch", function(event) {
    // Cache-First Strategy
    event.respondWith(
        caches
            .match(event.request)  // check if the request has already been cached
            .then(function(cached) {
                return cached || fetch(event.request) // otherwise request network
            })
    );
});
