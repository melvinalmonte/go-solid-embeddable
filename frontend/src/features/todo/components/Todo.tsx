import { Component, createSignal, For, Match, Switch } from "solid-js";
import { createTodo, deleteTodo, getTodos, updateTodo } from "../api";

const Todo: Component = () => {
  const [newTodo, setNewTodo] = createSignal("");

  // CRUD operations
  const createEntry = createTodo();
  const readTodos = getTodos();
  const updateEntry = updateTodo();
  const deleteEntry = deleteTodo();

  const onAdd = () => {
    if (newTodo() !== "") {
      createEntry.mutate(newTodo());
      setNewTodo("");
    }
  };

  return (
    <>
      <div class="overflow-x-auto p-4">
        <table class="table w-full">
          <thead>
            <tr>
              <th>Done</th>
              <th>Task</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            <Switch>
              <Match when={readTodos.isLoading} keyed>
                <p>Loading...</p>
              </Match>
              <Match when={readTodos.isError} keyed>
                <p>Error: {readTodos?.error?.message}</p>
              </Match>
              <Match when={readTodos.isSuccess} keyed>
                <For each={readTodos.data}>
                  {todo => (
                    <tr>
                      <th>
                        <input
                          type="checkbox"
                          checked={todo.done}
                          class="checkbox checkbox-primary"
                          onChange={() => updateEntry.mutate({ ...todo, done: !todo.done })}
                        />
                      </th>
                      <td>{todo.text}</td>
                      <td>
                        <button
                          onClick={() => deleteEntry.mutate({ ...todo, done: !todo.done })}
                          class="btn btn-outline btn-error"
                        >
                          Remove
                        </button>
                      </td>
                    </tr>
                  )}
                </For>
              </Match>
            </Switch>
          </tbody>
        </table>
        <div class="form-control p-4">
          <div class="input-group">
            <input
              type="text"
              placeholder="Add todo"
              class="input input-bordered input-primary"
              value={newTodo()}
              onInput={e => setNewTodo(e.currentTarget.value)}
            />
            <button class="btn btn-primary" onClick={() => onAdd()}>
              Add
            </button>
          </div>
        </div>
      </div>
    </>
  );
};

export default Todo;
