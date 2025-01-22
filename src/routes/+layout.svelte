<script>
	/**
	 * @typedef {Object} Props
	 * @property {import('svelte').Snippet} [children]
	 */

	/** @type {Props} */
	let { children } = $props();

	import '../app.css';
	children;
</script>

{@render children?.()}

<style>
	:global(body) {
		margin: 0;
		padding: 0;
	}
</style>
