<script>
 import { _ } from 'svelte-i18n';
</script>

<div class="separator">
    <kbd>{$_('purchased_grocery_separator')}</kbd>
</div>

<style>
 kbd {
  background-color: var(--primary);
 }
.separator {
  display: flex;
  flex-direction: column;
  align-items: center;
}
</style>
