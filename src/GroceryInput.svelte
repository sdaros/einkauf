<script>
 import { createEventDispatcher } from 'svelte';
 import { _ } from 'svelte-i18n';

 const dispatch = createEventDispatcher();
 let grocery = {
   id: 0,
   title: '',
   purchased: false
 };

 function addGrocery() {
     grocery.id = unixEpoch();
     grocery.title = grocery.title.trim();
     dispatch('message', {
         value: grocery
     });
     reset();
 }
 function reset() {
     grocery.title = '';
     grocery.id = 0;
 }
 function unixEpoch() {
     return new Date().getTime()
 }
</script>

<input
  autocomplete="off"
  autocorrect="off"
  autocapitalize="on"
  spellcheck="false"
  placeholder="{$_('grocery_entry_placeholder')}"
  bind:value={grocery.title}
  on:blur={addGrocery}
  on:keyup={e => {if (e.key === 'Enter') { addGrocery() }}}
/>
