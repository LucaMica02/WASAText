<script>
export default {};
</script>

<template>
  <div class="search-bar">
    <input
      type="text"
      v-model="searchQuery"
      placeholder="Search a username"
      class="search-bar-input"
    />
  </div>
</template>

<style>
.search-bar {
  display: inline;
  margin: 10px;
}

.search-bar-input {
  height: 20px;
  width: 200px;
}
</style>
