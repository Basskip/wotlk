import { Tooltip } from 'bootstrap';
import { title } from 'process';
import { Component } from './component.js';

export interface ContentBlockHeaderConfig {
  title: string,
  classes?: Array<string>,
  titleTag?: string,
  tooltip?: string,
}

export interface ContentBlockConfig {
	bodyClasses?: Array<string>,
  classes?: Array<string>,
  rootElem?: HTMLElement,
  header?: ContentBlockHeaderConfig,
}

export class ContentBlock extends Component {
  readonly headerElement: HTMLElement|null;
  readonly bodyElement: HTMLElement;

  readonly config: ContentBlockConfig;

	constructor(parent: HTMLElement, cssClass: string, config: ContentBlockConfig) {
		super(parent, 'content-block', config.rootElem);
    this.config = config;
		this.rootElem.classList.add(cssClass);

		if (config.classes) {
			this.rootElem.classList.add(...config.classes);
    }

    this.headerElement = this.buildHeader();
    this.bodyElement = this.buildBody();
	}

  private buildHeader(): HTMLElement|null {
    if (this.config.header && Object.keys(this.config.header).length) {
      let titleTag = this.config.header.titleTag || 'h6';
      let headerFragment = document.createElement('fragment');
      headerFragment.innerHTML = `
        <div class="content-block-header">
          <${titleTag}
            class="content-block-title"
            ${this.config.header.tooltip ? 'data-bs-toggle="tooltip"' : ''}
            ${this.config.header.tooltip ? `data-bs-title="${this.config.header.tooltip}"` : ''}
          >${this.config.header.title}</${titleTag}>
        </div>
      `;

      let header = headerFragment.children[0] as HTMLElement;
      
      if (this.config.header.classes) {
        header.classList.add(...this.config.header.classes);
      }

      if (this.config.header.tooltip)
        Tooltip.getOrCreateInstance(header.querySelector('.content-block-title') as HTMLElement);

      this.rootElem.appendChild(header);

      return header;
    } else {
      return null;
    }
  }

  private buildBody(): HTMLElement {
    let bodyElem = document.createElement('div');
    bodyElem.classList.add('content-block-body');

    this.rootElem.appendChild(bodyElem);

    return bodyElem;
  }
}
