<template>
    <div class="container momo-search">
        <slot></slot>
    </div>
</template>

<script>

export default {
  created: function () {
    // Initialize the entire application
    this.$store.dispatch('loadResults')
  }
}
</script>
