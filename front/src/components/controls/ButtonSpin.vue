<template>
  <b-button :type="type" size="lg" :variant="variant" :disabled="spinning || disabled">
    {{label}} <i v-if="spinning" class="fa fa-spinner fa-spin" aria-hidden="true"></i>
  </b-button>
</template>
<script>
export default {
  props: {
    label: {
      type: String,
      required: true
    },
    variant: {
      type: String,
      default: 'primary',
      validator: function (value) {
        return ['primary', 'secondary', 'success', 'warning', 'danger', 'link'].indexOf(value) !== -1
      }
    },
    spinning: {
      default: false,
      type: Boolean
    },
    disabled: {
      default: false,
      type: Boolean
    },
    type: {
      type: String,
      default: 'button',
      validator: function (value) {
        return ['button', 'submit'].indexOf(value) !== -1
      }
    },
  }
}
</script>
<style scoped>
</style>
