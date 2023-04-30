import type { Component } from "solid-js";
import { Todos } from "./features/todo";
import { Navigate, Route, Router, Routes } from "@solidjs/router";

const App: Component = () => {
  return (
    <Router>
      <Routes>
        <Route path="/home" component={Todos} />
        <Route path="*" element={<Navigate href="/home" />} />
      </Routes>
    </Router>
  );
};

export default App;
