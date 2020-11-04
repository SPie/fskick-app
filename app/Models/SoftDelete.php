<?php

namespace App\Models;

use Doctrine\ORM\Mapping as ORM;
use Gedmo\Mapping\Annotation as Gedmo;

/**
 * Trait SoftDelete
 *
 * @package App\Models
 */
trait SoftDelete
{
    /**
     * @ORM\Column(name="deleted_at", type="datetime", nullable=true)
     *
     * @var \DateTime|null
     */
    private ?\DateTime $deletedAt;

    /**
     * @param \DateTime|null $deletedAt
     *
     * @return $this
     */
    public function setDeletedAt(?\DateTime $deletedAt): self
    {
        $this->deletedAt = $deletedAt;

        return $this;
    }

    /**
     * @return \DateTime|null
     */
    public function getDeletedAt(): ?\DateTime
    {
        return $this->deletedAt;
    }

    /**
     * @return $this
     */
    public function restore()
    {
        $this->deletedAt = null;

        return $this;
    }

    /**
     * @return bool
     */
    public function isDeleted(): bool
    {
        return ($this->getDeletedAt() && (new \DateTime()) >= $this->getDeletedAt());
    }
}
