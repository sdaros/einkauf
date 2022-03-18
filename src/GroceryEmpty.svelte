<script>
 import { _ } from 'svelte-i18n';
</script>

<aside class="empty-groceries">
  <h5>{$_('empty_groceries_heading')}</h5>
  <p>{$_('empty_groceries_paragraph')}</p>
</aside>

<style>
 .empty-groceries {
  display: flex;
  flex-direction: column;
  align-items: center;
 }
</style>
