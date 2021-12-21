import { computed, ref } from 'vue';

type Favorite = {
  id: string;
  type: string;
  name: string;
};

export function useFavorites() {
  const LS_FAVORITES_KEY = 'favoriteStops';

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

  function addFavorite(favorite: Favorite) {
    favorites.value = [...favorites.value, favorite];
  }

  function removeFavorite(favorite: Favorite) {
    favorites.value = favorites.value.filter((f) => f.id !== favorite.id);
  }

  function isFavorite(id: string) {
    return favorites.value.some((f) => f.id === id);
  }

  return { favorites, addFavorite, removeFavorite, isFavorite };
}
