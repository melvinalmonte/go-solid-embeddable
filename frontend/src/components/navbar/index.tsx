import { Component, Show } from "solid-js";
import { A, useLocation } from "@solidjs/router";

const Navbar: Component = () => {
  const location = useLocation();

  const isActive = (path: string) => location.pathname === path;
  return (
    <div class="navbar bg-primary text-primary-content">
      <div class="container mx-auto">
        <a class=" normal-case text-xl">Another TODO App</a>
      </div>
    </div>
  );
};

export default Navbar;
