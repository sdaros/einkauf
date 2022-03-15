<script>
 import { createEventDispatcher } from 'svelte';
 export let item;

</script>
<label>
  <input
    type="checkbox"
    bind:checked={item.purchased}
  />
    { item.title }
</label>
<style>
</style>
