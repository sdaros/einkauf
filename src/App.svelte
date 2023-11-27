<script>
 import { svelteStore } from "./store.js";
 import { _ } from 'svelte-i18n';
 import GroceryInput from './GroceryInput.svelte';
 import GroceryItem from './GroceryItem.svelte';
 import GroceryEmpty from './GroceryEmpty.svelte';
 import PurchasedSeparator from './PurchasedSeparator.svelte';

 const logoNameFromFragment = (fragment) => {
     let logoName = fragment.split("_")[0];
     if (logoName !== "food") {
         logoName = "other"
     }
     return logoName
 }

 function dynamicManifest(fragment) {
     const logoName = logoNameFromFragment(fragment)
     const manifest = {
         name: "Einkauf App",
         short_name: "Einkauf",
         theme_color: "#2df321",
         background_color: "#fdfdfd",
         display: "fullscreen",
         orientation: "portrait",
         scope: "/",
         start_url: "/#" + fragment,
         icons: [
             {
                 src: logoName + "-64x64.png",
                 sizes: "64x64",
                 type: 'image/png'
             },
             {
                 src: logoName + "-192x192.png",
                 sizes: "192x192",
                 type: 'image/png'
             },
             {
                 src: logoName + "-512x512.png",
                 sizes: "512x512",
                 type: 'image/png',
                 purpose: 'any'
             },
             {
                 src: logoName + "-maskable-512x512.png",
                 sizes: "512x512"
                 type: 'image/png',
                 purpose: 'maskable'
             }
         ]
     }
     const stringManifest = JSON.stringify(manifest)
     const blob = new Blob([stringManifest], {type: 'application/json'})
     return URL.createObjectURL(blob)
 }
 const addGrocery = (event) => {
     const value = event.detail.value;
     if (!value.title) {
         return;
     }
     $svelteStore.groceries.push({...value});
 };
 // to force it
 const updatePurchasedState = (event, grocery) => {
     const value = event.detail.value;
     let updatedGrocery = {...grocery.title, ...grocery.id, purchased: value}
     $svelteStore.groceries.splice(indexAt(grocery), 1, updatedGrocery)
 };
 const clearGroceries = () => {
     const groceries = $svelteStore.groceries.filter(n => n.purchased);
     for (let grocery of groceries) {
         let idx = $svelteStore.groceries.indexOf(grocery)
         $svelteStore.groceries.splice(idx, 1)
     }
 };
 const byTitle = (a,b) => {
     let aTitle = a.title
     let bTitle = b.title
     if (aTitle < bTitle) {
         return -1
     }
     if (aTitle > bTitle) {
         return 1
     }
     return 0
 };
 $: unpurchasedGroceries = $svelteStore.groceries.filter(n => !n.purchased).sort(byTitle)
 $: purchasedGroceries = $svelteStore.groceries.filter(n => n.purchased).sort(byTitle)
 $: emptyGroceries = (unpurchasedGroceries.length === 0 && purchasedGroceries.length === 0)

</script>

<svelte:head>
<link rel="manifest" href={dynamicManifest(window.location.hash.substring(1))}>
<link rel="icon" href={"/assets/" + logoNameFromFragment(window.location.hash.substring(1)) + "-favicon.ico"} sizes="any">
<link rel="icon" href={"/assets/" + logoNameFromFragment(window.location.hash.substring(1)) + "-favicon.svg"} type="image/svg+xml">
<link rel="apple-touch-icon" href={"/assets/" + logoNameFromFragment(window.location.hash.substring(1)) + "-180x180.png"}>
</svelte:head>
<main class="container">
    <article>
        <header>
            <GroceryInput on:message={addGrocery}/>
        </header>
        {#if emptyGroceries}
            <GroceryEmpty />
        {/if}
        {#if unpurchasedGroceries.length > 0}
            <section>
                <ul>
                    {#each unpurchasedGroceries as grocery (grocery.id)}
                        <li>
                            <GroceryItem item={grocery} />
                        </li>
                    {/each}
                </ul>
            </section>
        {/if}
        {#if purchasedGroceries.length > 0}
            <PurchasedSeparator />
            <section>
                <ul>
                    {#each purchasedGroceries as grocery (grocery.id)}
                        <li>
                            <GroceryItem item={grocery} />
                        </li>
                    {/each}
                </ul>
            </section>
            <footer>
                <button
                    on:click={clearGroceries}
                    class="warning"
                >
                    {$_('delete_purchased_groceries')}
                </button>
            </footer>
        {/if}
        <footer/>
    </article>
</main>

<style>
 ul li {
     list-style: none;
 }
 .warning {
     background-color: var(--del-color);
     border-color: var(--del-color)
 }
</style>
