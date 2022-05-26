import { computed, ref } from 'vue';

type Favorite = {
  id: string;
  type: string;
  name: string;
};
const LS_FAVORITES_KEY = 'kiel-live-favorites-v1';

// migrate legacy favorites
type LegacyFavorite = {
  id: string;
  name: string;
  favorite: true;
};
const LS_LEGACY_FAVORITES_KEY = 'favoriteStops';
const legacyLocalStorageItem = localStorage.getItem(LS_LEGACY_FAVORITES_KEY);
if (legacyLocalStorageItem !== null) {
  const legacyFavorites = JSON.parse(legacyLocalStorageItem) as LegacyFavorite[];
  localStorage.setItem(
    LS_FAVORITES_KEY,
    JSON.stringify(legacyFavorites.map((f) => ({ id: `kvg-${f.id}`, name: f.name, type: 'bus-stop' }))),
  );
  localStorage.removeItem(LS_LEGACY_FAVORITES_KEY);
}

const favoritesRaw = ref<Favorite[]>(JSON.parse(localStorage.getItem(LS_FAVORITES_KEY) || '[]') as Favorite[]);

const favorites = computed({
  get() {
    return favoritesRaw.value;
  },
  set(_favorites: Favorite[]) {
    favoritesRaw.value = _favorites;
    localStorage.setItem(LS_FAVORITES_KEY, JSON.stringify(_favorites));
  },
});

function addFavorite({ id, name, type }: Favorite) {
  favorites.value = [...favorites.value, { id, name, type }];
}

function removeFavorite(favorite: Pick<Favorite, 'id'>) {
  favorites.value = favorites.value.filter((f) => f.id !== favorite.id);
}

function isFavorite(favorite: Pick<Favorite, 'id'>) {
  return favorites.value.some((f) => f.id === favorite.id);
}

export function useFavorites() {
  return { favorites, addFavorite, removeFavorite, isFavorite };
}
