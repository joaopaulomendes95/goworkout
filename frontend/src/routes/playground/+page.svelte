<script lang="ts">
	// some stuff from tutorial
	let header = 'playground to test stuff';

  // runes
  let count = $state(0);
  let numbers = $state([1, 2, 3, 4]);

  // derived runes
  let total = $derived(numbers.reduce((t, n) => t + n, 0));

  // effects
  let elapsed = $state(0);
  let interval = $state(1000);

  $effect(() => {
    const id = setInterval (() => {
      elapsed +1;
    }, interval);

    return() => {
      clearInterval(id);
    };
  });

  function increment() {
    count += 1;
  }

  function addNumber() {
    numbers.push(numbers.length + 1);
    $inspect(numbers).with(console.trace);
  }
</script>


<h1 class="text-2xl text-orange-300">{@html header}</h1>


<button onclick={increment}>
  Clicked {count}
  {count === 1 ? 'time' : 'times'}
</button>


<p>{numbers.join(' + ')} = {total}</p>


<button onclick={addNumber}>
  Add a number
</button>


<button onclick={() => interval /= 2}>speed up</button>
<button onclick={() => interval *= 2}>slow down</button>

<p>elapsed: {elapsed}</p>
