import { syncedStore, getYjsValue } from "@syncedstore/core";
import { svelteSyncedStore } from "@syncedstore/svelte";
import { IndexeddbPersistence } from "y-indexeddb";
import { WebsocketProvider } from "y-websocket";

// Create your SyncedStore store
export const store = syncedStore({ groceries: [], fragment: "xml" });
// Create Svelte Store for use in your components.
// You can treat this like any other store, including `bind`.
export const svelteStore = svelteSyncedStore(store);
const doc = getYjsValue(store);
const storeId = window.location.hash.slice(1);
const provider = new IndexeddbPersistence("einkauf_v2.0.0_" + storeId, doc);
const wsProvider = new WebsocketProvider("wss://einkauf.cip.li/ws", storeId, doc);
