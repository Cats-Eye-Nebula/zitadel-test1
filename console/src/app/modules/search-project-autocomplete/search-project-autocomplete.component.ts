import { COMMA, ENTER } from '@angular/cdk/keycodes';
import { Component, ElementRef, EventEmitter, Input, Output, ViewChild } from '@angular/core';
import { FormControl } from '@angular/forms';
import { MatAutocomplete, MatAutocompleteSelectedEvent } from '@angular/material/autocomplete';
import { MatChipInputEvent } from '@angular/material/chips';
import { from, merge } from 'rxjs';
import { debounceTime, switchMap, tap } from 'rxjs/operators';
import {
    ProjectGrantView,
    ProjectSearchKey,
    ProjectSearchQuery,
    ProjectView,
    SearchMethod,
} from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';

@Component({
    selector: 'app-search-project-autocomplete',
    templateUrl: './search-project-autocomplete.component.html',
    styleUrls: ['./search-project-autocomplete.component.scss'],
})
export class SearchProjectAutocompleteComponent {
    public selectable: boolean = true;
    public removable: boolean = true;
    public addOnBlur: boolean = true;
    public separatorKeysCodes: number[] = [ENTER, COMMA];
    public myControl: FormControl = new FormControl();
    public names: string[] = [];
    public projects: Array<ProjectGrantView.AsObject | ProjectView.AsObject | any> = [];
    public filteredProjects: Array<ProjectGrantView.AsObject | ProjectView.AsObject | any> = [];
    public isLoading: boolean = false;
    @ViewChild('nameInput') public nameInput!: ElementRef<HTMLInputElement>;
    @ViewChild('auto') public matAutocomplete!: MatAutocomplete;
    @Input() public singleOutput: boolean = false;
    @Output() public selectionChanged: EventEmitter<
        ProjectGrantView.AsObject[]
        | ProjectGrantView.AsObject
        | ProjectView.AsObject
        | ProjectView.AsObject[]
    > = new EventEmitter();
    constructor(private projectService: ProjectService) {
        this.myControl.valueChanges
            .pipe(
                debounceTime(200),
                tap(() => this.isLoading = true),
                switchMap(value => {
                    const query = new ProjectSearchQuery();
                    query.setKey(ProjectSearchKey.PROJECTSEARCHKEY_PROJECT_NAME);
                    query.setValue(value);
                    query.setMethod(SearchMethod.SEARCHMETHOD_CONTAINS_IGNORE_CASE);
                    return merge(
                        from(this.projectService.SearchGrantedProjects(10, 0, [query])),
                        from(this.projectService.SearchProjects(10, 0, [query])),
                    );
                }),
                // finalize(() => this.isLoading = false),
            ).subscribe((projects) => {
                this.isLoading = false;
                this.filteredProjects = projects.toObject().resultList;
                console.log(this.filteredProjects);
            });
    }

    public displayFn(project?: any): string | undefined {
        return (project && project.projectName) ? `${project.projectName}` :
            (project && project.name) ? `${project.name}` : undefined;
    }

    public add(event: MatChipInputEvent): void {
        if (!this.matAutocomplete.isOpen) {
            const input = event.input;
            const value = event.value;

            if ((value || '').trim()) {
                const index = this.filteredProjects.findIndex((project) => {
                    if (project?.projectName) {
                        return project.projectName === value;
                    } else if (project?.name) {
                        return project.name === value;
                    }
                });
                if (index > -1) {
                    if (this.projects && this.projects.length > 0) {
                        this.projects.push(this.filteredProjects[index]);
                    } else {
                        this.projects = [this.filteredProjects[index]];
                    }
                }
            }

            if (input) {
                input.value = '';
            }
        }
    }

    public remove(project: ProjectGrantView.AsObject): void {
        const index = this.projects.indexOf(project);

        if (index >= 0) {
            this.projects.splice(index, 1);
        }
    }

    public selected(event: MatAutocompleteSelectedEvent): void {
        console.log(event.option.value);
        if (this.singleOutput) {
            this.selectionChanged.emit(event.option.value);
        } else {
            if (this.projects && this.projects.length > 0) {
                this.projects.push(event.option.value);
            } else {
                this.projects = [event.option.value];
            }
            this.selectionChanged.emit(this.projects);

            this.nameInput.nativeElement.value = '';
            this.myControl.setValue(null);
        }
    }
}
