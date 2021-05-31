<template>
    <div class="tagger">
        <b-form @submit.prevent="onSubmit">
            <b-form-group
                id="input-group-1"
                label="Tag this object:"
                label-for="inline-form-input-tag"
                description="Add a single tag."
            >

            <b-input-group>
                <b-form-input
                    id="inline-form-input-tag"
                    v-model="tag"
                    class="mb-2 mr-sm-2 mb-sm-0"
                ></b-form-input>
                <b-input-group-append>
                    <b-button variant="primary" @click.prevent="addTag">Add tag</b-button>
                </b-input-group-append>
            </b-input-group>

            <ul id="tagger-tags" class="list-inline mt-3 mb-0">
                <li
                v-for="tag in form.tags"
                v-bind:id="tag"
                :key="tag"
                class="list-inline-item mb-3"
                >
                    <span class="bg-dark text-white rounded p-2">{{tag}} <a href="#" class="text-white" @click.prevent="removeTag(tag)">&times;</a></span>
                </li>
            </ul>
            </b-form-group>
            <b-button type="submit" variant="primary">Save</b-button>
        </b-form>
    </div>
</template>


<script>

import { mapState } from 'vuex'

export default {
  data() {
    return {
      tag: ''
    }
  },
  computed: {
    ...mapState([
      'form'
    ])
  },
  methods: {
    addTag() {
      this.$store.dispatch('addTag', { tag: this.tag })
      this.tag = ""
    },
    removeTag(tag) {
      this.$store.dispatch('removeTag', { tag: tag })
    },
  }
}
</script>