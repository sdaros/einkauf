<script>
 import { svelteStore } from "./store.js";
 import GroceryInput from './GroceryInput.svelte';
 import GroceryItem from './GroceryItem.svelte';
 import GroceryEmpty from './GroceryEmpty.svelte';
 import PurchasedSeparator from './PurchasedSeparator.svelte';

 const addGrocery = (event) => {
   const value = event.detail.value;
   if (!value.title) {
     return;
   }
   $svelteStore.groceries.push({...value});
 };
 const updatePurchasedState = (event, grocery) => {
	 const value = event.detail.value;
	 let updatedGrocery = {...grocery.title, ...grocery.id, purchased: value}
	 $svelteStore.groceries.splice(indexAt(grocery), 1, updatedGrocery)
 };
 const clearGroceries = () => {
      $svelteStore.groceries.splice(0, $svelteStore.groceries.length);
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

</script>

<main class="container">
{#if $svelteStore.groceries.length > 0}
<article>
 <header>
   <GroceryInput on:message={addGrocery}/>
 </header>
 <section>
  <ul>
   {#each $svelteStore.groceries.filter(n => !n.purchased).sort(byTitle) as grocery (grocery.id)}
     <li>
       <GroceryItem item={grocery} />
     </li>
   {/each}
 </ul>
 </section>
 <PurchasedSeparator />
 <section>
 <ul>
 {#each $svelteStore.groceries.filter(n => n.purchased).sort(byTitle) as grocery (grocery.id)}
 	<li>
 	  <GroceryItem item={grocery} />
 	</li>
 {/each}
 </ul>
 </section>
<footer>
  <button
	on:click|once={clearGroceries}
    class="warning"
  >
    Lösche alle gekaufte Zutaten
  </button>
</footer>
</article>
{:else}
<article>
	<header>
	<GroceryInput on:message={addGrocery}/>
	</header>
	<section>
		<GroceryEmpty />
	</section>
	<footer />
</article>
{/if}
</main>

<style>
ul li {
  list-style:none;
}
.warning {
	background-color: var(--del-color);
	border-color: var(--del-color)
 }
</style>
