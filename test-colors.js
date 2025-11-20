import blessed from 'blessed';

const screen = blessed.screen({
  mouse: true,
  keyboard: true,
  trueColor: true,
  smartCSR: true,
});

const box = blessed.box({
  parent: screen,
  top: 0,
  left: 0,
  right: 0,
  height: 15,
  padding: 1,
  style: {
    bg: 'white',
    fg: 'black',
  },
  tags: true,
  border: 'line',
});

const content = `{bold}Colors Test{/bold}
{red}This is red{/red}
{green}This is green{/green}
{blue}This is blue{/blue}
{yellow}This is yellow{/yellow}
{cyan}This is cyan{/cyan}
{magenta}This is magenta{/magenta}
{bold}This is bold{/bold}
Plain text`;

box.setContent(content);

screen.key(['q'], () => process.exit(0));
screen.render();
