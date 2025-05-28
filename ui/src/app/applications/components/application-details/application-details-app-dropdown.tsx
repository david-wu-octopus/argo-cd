import {DropDown} from 'argo-ui';
import * as React from 'react';

import {Context} from '../../../shared/context';
import {getAppUrl} from '../utils';
import {ApplicationList} from '../../../shared/models';

import './application-details-app-dropdown.scss';

export const ApplicationsDetailsAppDropdown = (props: {appName: string; opened: boolean; setOpened: (opened: boolean) => void; apps: ApplicationList}) => {
    const [appFilter, setAppFilter] = React.useState('');
    const [selectedIndex, setSelectedIndex] = React.useState(0);
    const ctx = React.useContext(Context);
    const selectedRef = React.useRef<HTMLLIElement>(null);

    // Scroll selected item into view when selection changes
    React.useEffect(() => {
        if (selectedRef.current) {
            selectedRef.current.scrollIntoView({
                block: 'nearest',
                behavior: 'smooth'
            });
        }
    }, [selectedIndex]);

    const handleKeyDown = (e: React.KeyboardEvent, filteredApps: any[]) => {
        switch (e.key) {
            case 'ArrowDown':
                e.preventDefault();
                if (selectedIndex < filteredApps.length - 1) {
                    const nextIndex = selectedIndex + 1;
                    setSelectedIndex(nextIndex);
                }
                break;
            case 'ArrowUp':
                e.preventDefault();
                if (selectedIndex > 0) {
                    const prevIndex = selectedIndex - 1;
                    setSelectedIndex(prevIndex);
                }
                break;
            case 'Enter':
                e.preventDefault();
                e.stopPropagation();
                if (filteredApps[selectedIndex]) {
                    ctx.navigation.goto(`/${getAppUrl(filteredApps[selectedIndex])}`);
                }
                break;
        }
    };

    const filteredApps = props.apps.items
        .filter(app => {
            return appFilter.length === 0 || app.metadata.name.toLowerCase().includes(appFilter.toLowerCase());
        })
        .slice(0, 100);

    return (
        <DropDown
            onOpenStateChange={props.setOpened}
            isMenu={true}
            anchor={() => (
                <>
                    <i className='fa fa-search' /> <span>{props.appName}</span>
                </>
            )}>
            {props.opened && (
                <ul className='application-details-app-dropdown'>
                    <li className='application-details-app-dropdown__search'>
                        <input
                            className='argo-field'
                            value={appFilter}
                            onChange={e => setAppFilter(e.target.value)}
                            onKeyDown={e => handleKeyDown(e, filteredApps)}
                            ref={el =>
                                el &&
                                setTimeout(() => {
                                    if (el) {
                                        el.focus();
                                    }
                                }, 100)
                            }
                        />
                    </li>
                    <div className='application-details-app-dropdown__list'>
                        {filteredApps.map((app, index) => (
                            <li
                                key={app.metadata.name}
                                onClick={() => {
                                    ctx.navigation.goto(`/${getAppUrl(app)}`);
                                }}
                                className={`application-details-app-dropdown__item ${index === selectedIndex ? 'selected' : ''}`}
                                onMouseEnter={() => setSelectedIndex(index)}
                                ref={index === selectedIndex ? selectedRef : null}>
                                {app.metadata.name} {app.metadata.name === props.appName && ' (current)'}
                            </li>
                        ))}
                    </div>
                </ul>
            )}
        </DropDown>
    );
};
