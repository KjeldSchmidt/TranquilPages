import { Component, EventEmitter, Output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { BookService, Book } from '../../services/book.service';

@Component({
  selector: 'app-add-book-modal',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="modal-overlay" *ngIf="isOpen" (click)="close()">
      <div class="modal-content" (click)="$event.stopPropagation()">
        <h2>Add New Book</h2>
        <form (ngSubmit)="onSubmit()">
          <div class="form-group">
            <label for="title">Title</label>
            <input type="text" id="title" [(ngModel)]="book.title" name="title" required>
          </div>
          <div class="form-group">
            <label for="author">Author</label>
            <input type="text" id="author" [(ngModel)]="book.author" name="author" required>
          </div>
          <div class="form-group">
            <label for="rating">Rating</label>
            <input type="number" id="rating" [(ngModel)]="book.rating" name="rating" min="0" max="5" required>
          </div>
          <div class="form-group">
            <label for="comment">Comment</label>
            <textarea id="comment" [(ngModel)]="book.comment" name="comment"></textarea>
          </div>
          <div class="button-group">
            <button type="button" (click)="close()">Cancel</button>
            <button type="submit">Add Book</button>
          </div>
        </form>
      </div>
    </div>
  `,
  styles: [`
    .modal-overlay {
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background-color: rgba(0, 0, 0, 0.5);
      display: flex;
      justify-content: center;
      align-items: center;
      z-index: 1000;
    }

    .modal-content {
      background: white;
      padding: 2rem;
      border-radius: 8px;
      width: 90%;
      max-width: 500px;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    }

    .form-group {
      margin-bottom: 1rem;
    }

    label {
      display: block;
      margin-bottom: 0.5rem;
      font-weight: 500;
    }

    input, textarea {
      width: 100%;
      padding: 0.5rem;
      border: 1px solid #ddd;
      border-radius: 4px;
      font-size: 1rem;
    }

    textarea {
      min-height: 100px;
      resize: vertical;
    }

    .button-group {
      display: flex;
      justify-content: flex-end;
      gap: 1rem;
      margin-top: 1.5rem;
    }

    button {
      padding: 0.5rem 1rem;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      font-size: 1rem;
      transition: background-color 0.2s;
    }

    button[type="button"] {
      background-color: #f0f0f0;
      color: #333;
    }

    button[type="submit"] {
      background-color: #4CAF50;
      color: white;
    }

    button:hover {
      opacity: 0.9;
    }
  `]
})
export class AddBookModalComponent {
  @Output() bookAdded = new EventEmitter<void>();
  @Output() modalClosed = new EventEmitter<void>();

  isOpen = false;
  book: Book = {
    title: '',
    author: '',
    rating: 0,
    comment: ''
  };

  constructor(private bookService: BookService) {}

  open() {
    this.isOpen = true;
  }

  close() {
    this.isOpen = false;
    this.resetForm();
    this.modalClosed.emit();
  }

  private resetForm() {
    this.book = {
      title: '',
      author: '',
      rating: 0,
      comment: ''
    };
  }

  onSubmit() {
    this.bookService.createBook(this.book).subscribe({
      next: (createdBook) => {
        this.book = createdBook; // Store the complete book data including _id
        this.bookAdded.emit();
        this.close();
      },
      error: (error) => {
        console.error('Failed to create book:', error);
        // TODO: Add error handling UI
      }
    });
  }
} 