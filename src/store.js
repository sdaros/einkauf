import { syncedStore, getYjsValue } from "@syncedstore/core";
import { svelteSyncedStore } from "@syncedstore/svelte";
import { IndexeddbPersistence } from "y-indexeddb";
import { WebsocketProvider } from "y-websocket";

// Create your SyncedStore store
export const store = syncedStore({ groceries: [], fragment: "xml" });
// Create Svelte Store for use in your components.
// You can treat this like any other store, including `bind`.
export const svelteStore = svelteSyncedStore(store);

// Create a document that syncs automatically using Y-WebRTC
const doc = getYjsValue(store);
const provider = new IndexeddbPersistence("my-document-id", doc);
const wsProvider = new WebsocketProvider("ws://192.168.0.108:8001", "einkauf", doc);
