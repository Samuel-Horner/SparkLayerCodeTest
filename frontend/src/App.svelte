<script lang="ts">
  import Todo from "./lib/Todo.svelte";
  import type { TodoItem } from "./lib/types";

  const backend_url = "http://localhost:8080/"

  let todos: TodoItem[] = $state([]);

  async function fetchTodos() {
    try {
      const response = await fetch(backend_url);

      if (response.status !== 200) {
        console.error("Error fetching data. Response status not 200");
        return;
      }

      todos = await response.json();

    } catch (e) {
      console.error("Could not connect to server. Ensure it is running.", e);
    }
  }

  function validateForm(formData: FormData) : boolean {
    if (formData.get("title") == "") return false;
    if (formData.get("description") == "") return false;

    return true;
  }

  async function submitTodo(e: SubmitEvent) {
    e.preventDefault(); // Cancel default url change behaviour

    // https://stackoverflow.com/questions/64527549/svelte-form-onsubmit-type-typescript
    const form = e.target as HTMLFormElement
    const formData = new FormData(form) 

    if (!validateForm(formData)) {
      alert("Invalid form values.");
      return;
    }

    try {
      const response = await fetch(backend_url, {
        method: "POST",
        body: JSON.stringify({
          title: formData.get("title"),
          description: formData.get("description")
        })
      });

      if (response.status !== 200) {
        console.error("Error submitting todo. Response status not 200");
        return;
      }

      todos.push(await response.json());
  
      form.reset();

    } catch (e) {
      console.error("Could not connect to server. Ensure it is running.", e);
    }
  }

  async function removeTodo(index: number) {
    if (index >= todos.length) {
      // Error not alert since this is only trigerable through console
      console.error("Invalid remove index.");
      return;
    }

    try {
      const response = await fetch(backend_url, {
        method: "DELETE",
        body: JSON.stringify({
          index: index
        })
      });

      if (response.status !== 200) {
        console.error("Error removing todo. Response status not 200");
        return;
      }

      todos.splice(index, 1);
    } catch (e) {
      console.error("Could not connect to server. Ensure it is running.", e);
    }
  }

  // Initially fetch todos on page load
  $effect(() => {
    fetchTodos();
  });
</script>

<main class="app">
  <header class="app-header">
    <h1>TODO</h1>
  </header>

  <div class="todo-list">
    {#each todos as todo, index}
      <Todo title={todo.title} description={todo.description} remove_function={() => removeTodo(index)} />
    {/each}
  </div>

  <h2 class="todo-list-form-header">Add a Todo</h2>
  <form class="todo-list-form" onsubmit={submitTodo}>
    <input placeholder="Title" name="title" />
    <input placeholder="Description" name="description" />
    <button>Add Todo</button>
  </form>
</main>

<style>
  .app {
    color: white;
    background-color: #282c34;

    text-align: center;
    font-size: 24px;

    min-height: 100vh;
    padding: 20px;
  }

  .app-header {
    font-size: calc(10px + 4vmin);
    margin-top: 50px;
  }

  .todo-list {
    margin: 50px 100px 0px 100px;
  }

  .todo-list-form-header {
    margin-top: 100px;
  }

  .todo-list-form {
    margin-top: 10px;
  }
</style>
