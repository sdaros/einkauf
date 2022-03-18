import { addMessages, init, getLocaleFromNavigator } from 'svelte-i18n';

import en from './lang/en.json';
import de from './lang/de.json';

addMessages('en', en);
addMessages('de', de);

init({
  fallbackLocale: 'en',
  initialLocale: getLocaleFromNavigator(),
});
