<template>
  <div class="delete-item">
    <i class="fa delete-item-i fa-times-circle" @click.stop="showModal" aria-hidden="true"></i>
    <b-modal ref="deleteModalRef" hide-footer title="Deletion">
      <div class="d-block text-center">
        <h3>Do you realy want to delete {{itemLabel}} ?</h3>
      </div>
      <button-spin 
           @click.native="hideModal" 
           id="cancel-deletion-button"
           type="button"
           variant="secondary"
           :disabled="isDeleting"
           label="Cancel"
           >
      </button-spin>
        <button-spin 
           @click.native="deleteItem" 
           id="do-deletion-button"
           type="button"
           variant="danger"
           :disabled="isDeleting"
           :spinning="isDeleting"
           label="Delete it !"
           >
        </button-spin>
    </b-modal>
  </div>
</template>
<script>
import ButtonSpin from '@/components/controls/ButtonSpin';
// todo tests
export default {
  components: {
    ButtonSpin,
  },
  props: {
    onDelete: {
      type: Function,
      required: true,
    },
    itemLabel: {
      type: String,
      required: true,
    }
  },
  data () {
    return {
      isDeleting: false,
    };
  },
  methods: {
    showModal () {
      this.$refs.deleteModalRef.show();
    },
    hideModal () {
      this.$refs.deleteModalRef.hide();
    },
    deleteItem: async function () {
      try {
        this.isDeleting = true;
        await this.onDelete();
				this.$emit('item-deleted')
        this.hideModal();
      } catch(err) {
        // todo show error
      } finally {
        this.isDeleting = false;
      }
    }
  },
}
</script>
<style scoped>
.delete-item {
  display: inline-block;
}
.delete-item-i:hover {
  color: var(--danger);
  cursor: pointer;
}
</style>
